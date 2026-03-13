package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project-mgmt/backend/internal/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBusinessLogic(t *testing.T) {
	// Seed data for the test
	SeedAll(testDB)

	t.Run("WBSProgressAndCostRollup", func(t *testing.T) {
		token := LoginAs(t, "admin@example.com", "password")
		projectIDStr := fmt.Sprintf("PRJ-ROLLUP-%d", time.Now().Unix())

		// 1. Create Project
		projReq := map[string]interface{}{
			"project_id":     projectIDStr,
			"project_name":   "Rollup Test Project",
			"project_status": "Active", // Correct field name
		}
		resp := MakeRequest(t, "POST", "/api/v1/projects", token, projReq)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var proj entity.Project
		json.Unmarshal(body, &proj)
		pID := proj.ID
		assert.NotEqual(t, 0, pID)

		// 2. Create Parent Phase (Dates: Today to +30 days)
		now := time.Now().UTC().Truncate(time.Hour)
		startDate := now.Format(time.RFC3339)
		endDate := now.AddDate(0, 0, 30).Format(time.RFC3339)

		phaseReq := map[string]interface{}{
			"project_id":         pID,
			"title":              "Phase 1",
			"type":               "Phase",
			"planned_start_date": startDate,
			"planned_end_date":   endDate,
		}
		resp = MakeRequest(t, "POST", fmt.Sprintf("/api/v1/projects/%d/wbs", pID), token, phaseReq)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body, _ = io.ReadAll(resp.Body)
		var phaseResp struct {
			Data entity.WBSNode `json:"data"`
		}
		json.Unmarshal(body, &phaseResp)
		parentPath := phaseResp.Data.Path
		parentID := phaseResp.Data.ID

		// 3. Create Child Task 1 (PV=100) (Dates: +1 to +10 days)
		task1Req := map[string]interface{}{
			"project_id":         pID,
			"title":              "Task 1.1",
			"type":               "Task",
			"parent_path":        parentPath,
			"planned_start_date": now.AddDate(0, 0, 1).Format(time.RFC3339),
			"planned_end_date":   now.AddDate(0, 0, 10).Format(time.RFC3339),
			"planned_value":      100,
		}
		resp = MakeRequest(t, "POST", fmt.Sprintf("/api/v1/projects/%d/wbs", pID), token, task1Req)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body, _ = io.ReadAll(resp.Body)
		var task1Resp struct {
			Data entity.WBSNode `json:"data"`
		}
		json.Unmarshal(body, &task1Resp)
		task1ID := task1Resp.Data.ID

		// 4. Create Child Task 2 (PV=300) (Dates: +11 to +20 days)
		task2Req := map[string]interface{}{
			"project_id":         pID,
			"title":              "Task 1.2",
			"type":               "Task",
			"parent_path":        parentPath,
			"planned_start_date": now.AddDate(0, 0, 11).Format(time.RFC3339),
			"planned_end_date":   now.AddDate(0, 0, 20).Format(time.RFC3339),
			"planned_value":      300,
		}
		resp = MakeRequest(t, "POST", fmt.Sprintf("/api/v1/projects/%d/wbs", pID), token, task2Req)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body, _ = io.ReadAll(resp.Body)
		var task2Resp struct {
			Data entity.WBSNode `json:"data"`
		}
		json.Unmarshal(body, &task2Resp)
		_ = task2Resp.Data.ID

		// 5. Update Task 1: Progress=100%, AC=50
		updateReq1 := map[string]interface{}{
			"id":            task1ID,
			"project_id":    pID,
			"title":         "Task 1.1 (Updated)",
			"type":          "Task",
			"progress":      100.0,
			"planned_value": 100.0,
			"actual_cost":   50.0,
		}
		resp = MakeRequest(t, "PUT", fmt.Sprintf("/api/v1/projects/%d/wbs/%d", pID, task1ID), token, updateReq1)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// 6. Verify Parent (Phase 1) Rollup
		// Expected Parent:
		// PV = 100 + 300 = 400
		// AC = 50 + 0 = 50
		// EV = (1.0 * 100) + (0 * 300) = 100
		// Progress = (100 / 400) * 100 = 25%
		resp = MakeRequest(t, "GET", fmt.Sprintf("/api/v1/projects/%d/wbs/%d", pID, parentID), token, nil)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ = io.ReadAll(resp.Body)
		var rolledUpParent struct {
			Data entity.WBSNode `json:"data"`
		}
		json.Unmarshal(body, &rolledUpParent)

		// Verification with some tolerance for float precision if needed
		assert.Equal(t, 400.0, rolledUpParent.Data.PlannedValue)
		assert.Equal(t, 50.0, rolledUpParent.Data.ActualCost)
		assert.Equal(t, 25.0, rolledUpParent.Data.Progress)

		// 7. Verify EVM Stats via API
		// SPI = EV / PV = 100 / 400 = 0.25
		// CPI = EV / AC = 100 / 50 = 2.0
		resp = MakeRequest(t, "GET", fmt.Sprintf("/api/v1/projects/%d/pmi-stats", pID), token, nil)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ = io.ReadAll(resp.Body)
		// Actually pmiStats data mapping is simple
		var pmiResult struct {
			Data map[string]interface{} `json:"data"`
		}
		json.Unmarshal(body, &pmiResult)

		assert.Equal(t, 400.0, pmiResult.Data["pv"])
		assert.Equal(t, 100.0, pmiResult.Data["ev"])
		assert.Equal(t, 50.0, pmiResult.Data["ac"])
		assert.Equal(t, 0.25, pmiResult.Data["spi"])
		assert.InDelta(t, 2.0, pmiResult.Data["cpi"], 0.001)
	})

	t.Run("SnapshotCaptureAndTrends", func(t *testing.T) {
		token := LoginAs(t, "admin@example.com", "password")

		// 1. Capture Snapshot
		resp := MakeRequest(t, "POST", "/api/v1/reporting/snapshots/capture", token, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to capture snapshot: Status %d, Body: %s", resp.StatusCode, string(body))
		}

		// 2. Verify Trends List
		resp = MakeRequest(t, "GET", "/api/v1/projects", token, nil)
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var projList struct {
			Data []entity.Project `json:"data"`
		}
		json.Unmarshal(body, &projList)

		// Find a project that has snapshots
		for _, p := range projList.Data {
			statusStr := "<nil>"
			if p.ProjectStatus != nil {
				statusStr = *p.ProjectStatus
			}
			fmt.Printf("DEBUG: Found project ID=%d, Name=%s, string-status=%s, id-status=%v\n", p.ID, p.ProjectName, statusStr, p.ProjectStatusID)
			resp = MakeRequest(t, "GET", fmt.Sprintf("/api/v1/reporting/projects/%d/trends", p.ID), token, nil)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, _ = io.ReadAll(resp.Body)
				fmt.Printf("DEBUG: Trends for project %d: %s\n", p.ID, string(body))
				var trends []entity.ProjectSnapshot
				json.Unmarshal(body, &trends)
				if len(trends) > 0 {
					return // PASS
				}
			} else {
				body, _ = io.ReadAll(resp.Body)
				fmt.Printf("DEBUG: Failed to get trends for project %d, stat %d: %s\n", p.ID, resp.StatusCode, string(body))
			}
		}
		t.Errorf("No project found with snapshots after capture. Total projects checked: %d", len(projList.Data))
	})

	t.Run("ResourceWorkloadCalculation", func(t *testing.T) {
		token := LoginAs(t, "admin@example.com", "password")

		// 1. Get Workload
		start := time.Now().Format("2006-01-02")
		end := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		urlPath := fmt.Sprintf("/api/v1/resources/workload?start_date=%s&end_date=%s", start, end)

		resp := MakeRequest(t, "GET", urlPath, token, nil)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var workload entity.WorkloadOverview
		json.Unmarshal(body, &workload)

		assert.Equal(t, start, workload.StartDate)
		assert.Equal(t, end, workload.EndDate)
	})
}
