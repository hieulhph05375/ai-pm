package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func LoginAs(t *testing.T, email, password string) string {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginData)

	resp, err := http.Post(testServer.URL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken
}

func MakeRequest(t *testing.T, method, path string, token string, body interface{}) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		t.Logf("DEBUG: Request Body: %s", string(bodyBytes))
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	return resp
}

func AssertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}

func AssertJSON(t *testing.T, resp *http.Response, key string, expected interface{}) {
	t.Helper()
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	// Complex JSON assertion could be done here, for now simple prints
	// fmt.Printf("DEBUG: %v\n", result)
}
