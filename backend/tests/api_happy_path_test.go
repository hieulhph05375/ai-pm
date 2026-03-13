package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-mgmt/backend/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStakeholderAPI_HappyPath(t *testing.T) {
	token := LoginAs(t, "admin@example.com", "password")

	t.Run("List Stakeholders", func(t *testing.T) {
		resp := MakeRequest(t, "GET", "/api/v1/stakeholders?page=1&limit=10", token, nil)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)

		var result struct {
			Data []entity.Stakeholder `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		// We just check status code to ensure it works
	})

	t.Run("Create and Delete Stakeholder", func(t *testing.T) {
		stakeholder := entity.Stakeholder{
			Name:         "API_TEST_Stakeholder",
			Role:         "Tester",
			Organization: "Test Corp",
			Email:        "api_test@example.com",
		}

		// Create
		resp := MakeRequest(t, "POST", "/api/v1/stakeholders", token, stakeholder)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusCreated)

		var created entity.Stakeholder
		json.NewDecoder(resp.Body).Decode(&created)
		assert.NotZero(t, created.ID)

		// Delete
		respDel := MakeRequest(t, "DELETE", fmt.Sprintf("/api/v1/stakeholders/%d", created.ID), token, nil)
		defer respDel.Body.Close()
		AssertStatus(t, respDel, http.StatusOK)
	})
}

func TestProjectAPI_HappyPath(t *testing.T) {
	token := LoginAs(t, "admin@example.com", "password")

	t.Run("List Projects", func(t *testing.T) {
		resp := MakeRequest(t, "GET", "/api/v1/projects?page=1&limit=10&search=&status=", token, nil)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("Create Project", func(t *testing.T) {
		project := entity.Project{
			ProjectID:   "API-TEST-999", // High ID to avoid collision
			ProjectName: "API Test Project Alpha",
		}

		resp := MakeRequest(t, "POST", "/api/v1/projects", token, project)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusCreated)
	})
}

func TestWBSAPI_HappyPath(t *testing.T) {
	token := LoginAs(t, "admin@example.com", "password")

	// Create a project first to ensure we have one
	project := entity.Project{
		ProjectID:   "WBS-TEST-PRJ",
		ProjectName: "WBS Test Project",
	}
	respPrj := MakeRequest(t, "POST", "/api/v1/projects", token, project)
	defer respPrj.Body.Close()
	AssertStatus(t, respPrj, http.StatusCreated)

	var createdPrj entity.Project
	json.NewDecoder(respPrj.Body).Decode(&createdPrj)

	t.Run("Get WBS Tree", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/projects/%d/wbs", createdPrj.ID)
		resp := MakeRequest(t, "GET", path, token, nil)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)
	})
}

func TestCategoryAPI_HappyPath(t *testing.T) {
	token := LoginAs(t, "admin@example.com", "password")

	t.Run("List Category Types", func(t *testing.T) {
		resp := MakeRequest(t, "GET", "/api/v1/category-types", token, nil)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("List Categories", func(t *testing.T) {
		resp := MakeRequest(t, "GET", "/api/v1/categories?page=1&limit=10", token, nil)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)
	})
}

func TestAuthAPI_HappyPath(t *testing.T) {
	t.Run("Login Happy Path", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "admin@example.com",
			"password": "password",
		}
		resp := MakeRequest(t, "POST", "/api/v1/auth/login", "", loginData)
		defer resp.Body.Close()
		AssertStatus(t, resp, http.StatusOK)

		var result struct {
			AccessToken string `json:"access_token"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.NotEmpty(t, result.AccessToken)
	})
}
