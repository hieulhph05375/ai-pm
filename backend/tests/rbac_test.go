package tests

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

type rbacCase struct {
	Method   string
	Path     string
	Body     interface{}
	Expected map[string]int // role -> expected status
}

func TestRBACMatrix(t *testing.T) {
	// 1. Prepare Roles
	roles := []string{"admin", "pmo", "projectmanager", "teamlead", "member", "viewer"}
	tokens := make(map[string]string)

	for _, role := range roles {
		email := role + "@example.com"
		tokens[role] = LoginAs(t, email, "password")
	}

	// 2. Define Matrix Cases
	// true = 200/201, false = 403
	expects := func(admin, pmo, pm, tl, mem, view bool) map[string]int {
		m := make(map[string]int)
		status := func(b bool) int {
			if b {
				return 200
			}
			return http.StatusForbidden
		}
		m["admin"] = status(admin)
		m["pmo"] = status(pmo)
		m["projectmanager"] = status(pm)
		m["teamlead"] = status(tl)
		m["member"] = status(mem)
		m["viewer"] = status(view)
		return m
	}

	cases := []rbacCase{
		{
			Method:   "GET",
			Path:     "/api/v1/projects",
			Expected: expects(true, true, true, true, true, true),
		},
		{
			Method:   "POST",
			Path:     "/api/v1/projects",
			Body:     map[string]interface{}{"project_id": "PRJ-RBAC-001", "project_name": "RBAC Project", "project_status": "Running", "planned_start_date": "2026-03-10T00:00:00Z", "planned_end_date": "2026-06-10T00:00:00Z"},
			Expected: expects(true, false, true, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/users",
			Expected: expects(true, false, false, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/portfolio/overview",
			Expected: expects(true, true, false, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/settings",
			Expected: expects(true, true, true, false, false, false),
		},
		{
			Method:   "POST",
			Path:     "/api/v1/reporting/snapshots/capture",
			Body:     map[string]interface{}{},
			Expected: expects(true, true, false, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/resources/workload",
			Expected: expects(true, true, true, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/roles",
			Expected: expects(true, false, false, false, false, false),
		},
		{
			Method:   "GET",
			Path:     "/api/v1/holidays",
			Expected: expects(true, true, true, true, true, true),
		},
		{
			Method:   "POST",
			Path:     "/api/v1/holidays",
			Body:     map[string]interface{}{"name": "RBAC Holiday", "date": "2027-12-25", "type": "State", "is_recurring": true},
			Expected: expects(true, false, false, false, false, false),
		},
	}

	// 3. Execution
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %s", tc.Method, tc.Path), func(t *testing.T) {
			for role, expectedStatus := range tc.Expected {
				t.Run(role, func(t *testing.T) {
					body := tc.Body
					if tc.Method == "POST" && tc.Path == "/api/v1/projects" {
						if m, ok := body.(map[string]interface{}); ok {
							// Deep copy and update project_id
							newBody := make(map[string]interface{})
							for k, v := range m {
								newBody[k] = v
							}
							newBody["project_id"] = fmt.Sprintf("PRJ-RBAC-%s", role)
							body = newBody
						}
					}

					resp := MakeRequest(t, tc.Method, tc.Path, tokens[role], body)
					defer resp.Body.Close()

					resBody, _ := io.ReadAll(resp.Body)

					if expectedStatus == 200 {
						// Allow 200 OK or 201 Created for successful operations
						if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
							t.Errorf("RBAC_RESULT: role=%s method=%s path=%s status=%d expected=SUCCESS (200/201) FAIL. Body: %s",
								role, tc.Method, tc.Path, resp.StatusCode, string(resBody))
						} else {
							fmt.Printf("RBAC_RESULT: role=%s method=%s path=%s status=%d expected=SUCCESS PASS\n",
								role, tc.Method, tc.Path, resp.StatusCode)
						}
					} else {
						if resp.StatusCode != expectedStatus {
							t.Errorf("RBAC_RESULT: role=%s method=%s path=%s status=%d expected=%d FAIL. Body: %s",
								role, tc.Method, tc.Path, resp.StatusCode, expectedStatus, string(resBody))
						} else {
							fmt.Printf("RBAC_RESULT: role=%s method=%s path=%s status=%d expected=%d PASS\n",
								role, tc.Method, tc.Path, resp.StatusCode, expectedStatus)
						}
					}
				})
			}
		})
	}
}
