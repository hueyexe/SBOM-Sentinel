package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/database"
	"github.com/chrisclapham/SBOM-Sentinel/internal/transport/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer represents a test server instance for integration testing
type TestServer struct {
	Server   *httptest.Server
	Database *database.SQLiteRepository
	DBPath   string
}

// SetupTestServer creates a new test server with temporary database for integration testing
func SetupTestServer(t *testing.T) *TestServer {
	// Create temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_sentinel.db")

	// Initialize database
	repo, err := database.NewSQLiteRepository(dbPath)
	require.NoError(t, err, "Failed to initialize test database")

	// Create HTTP server with real handlers
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"sbom-sentinel"}`))
	})
	
	// API v1 routes
	mux.HandleFunc("/api/v1/sboms", rest.SubmitSBOMHandler(repo))
	mux.HandleFunc("/api/v1/sboms/get", rest.GetSBOMHandler(repo))
	mux.HandleFunc("/api/v1/sboms/", rest.AnalyzeSBOMHandler(repo))

	// Create test server
	server := httptest.NewServer(mux)

	return &TestServer{
		Server:   server,
		Database: repo,
		DBPath:   dbPath,
	}
}

// Cleanup closes the test server and cleans up resources
func (ts *TestServer) Cleanup() {
	if ts.Server != nil {
		ts.Server.Close()
	}
	if ts.Database != nil {
		ts.Database.Close()
	}
	if ts.DBPath != "" {
		os.Remove(ts.DBPath)
	}
}

// SubmitSBOMResponse represents the response from SBOM submission
type SubmitSBOMResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// AnalysisResponse represents the response from SBOM analysis
type AnalysisResponse struct {
	SBOMID  string                `json:"sbom_id"`
	Results []map[string]interface{} `json:"results"`
	Summary map[string]interface{}   `json:"summary"`
}

func TestHealthEndpoint(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup()

	resp, err := http.Get(ts.Server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var healthResp map[string]string
	err = json.Unmarshal(body, &healthResp)
	require.NoError(t, err)

	assert.Equal(t, "ok", healthResp["status"])
	assert.Equal(t, "sbom-sentinel", healthResp["service"])
}

func TestCompleteAPIWorkflow(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup()

	// Step 1: Submit SBOM file
	t.Log("Step 1: Submitting SBOM file...")
	
	// Create test SBOM data
	testSBOM := createTestSBOM()
	sbomJSON, err := json.Marshal(testSBOM)
	require.NoError(t, err)

	// Create multipart request
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("sbom", "test-sbom.json")
	require.NoError(t, err)
	
	_, err = part.Write(sbomJSON)
	require.NoError(t, err)
	
	err = writer.Close()
	require.NoError(t, err)

	// Submit SBOM
	req, err := http.NewRequest("POST", ts.Server.URL+"/api/v1/sboms", &requestBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify submission response
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "SBOM submission should return 201 Created")

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var submitResp SubmitSBOMResponse
	err = json.Unmarshal(respBody, &submitResp)
	require.NoError(t, err)

	assert.NotEmpty(t, submitResp.ID, "Response should contain SBOM ID")
	assert.Equal(t, "SBOM submitted successfully", submitResp.Message)
	
	sbomID := submitResp.ID
	t.Logf("✓ SBOM submitted successfully with ID: %s", sbomID)

	// Step 2: Retrieve the submitted SBOM
	t.Log("Step 2: Retrieving submitted SBOM...")
	
	getURL := fmt.Sprintf("%s/api/v1/sboms/get?id=%s", ts.Server.URL, sbomID)
	resp, err = http.Get(getURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "SBOM retrieval should return 200 OK")

	respBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var retrievedSBOM map[string]interface{}
	err = json.Unmarshal(respBody, &retrievedSBOM)
	require.NoError(t, err)

	assert.Equal(t, sbomID, retrievedSBOM["id"])
	assert.Equal(t, "Test Application", retrievedSBOM["name"])
	
	t.Logf("✓ SBOM retrieved successfully")

	// Step 3: Analyze SBOM with license agent only (default)
	t.Log("Step 3: Analyzing SBOM with license agent...")
	
	analyzeURL := fmt.Sprintf("%s/api/v1/sboms/%s/analyze", ts.Server.URL, sbomID)
	req, err = http.NewRequest("POST", analyzeURL, nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "SBOM analysis should return 200 OK")

	respBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var analysisResp AnalysisResponse
	err = json.Unmarshal(respBody, &analysisResp)
	require.NoError(t, err)

	// Verify analysis response
	assert.Equal(t, sbomID, analysisResp.SBOMID)
	assert.NotNil(t, analysisResp.Results, "Analysis should return results")
	assert.NotNil(t, analysisResp.Summary, "Analysis should return summary")

	// Check that we have license findings (GPL-3.0-only should be detected)
	assert.Greater(t, len(analysisResp.Results), 0, "Should have at least one finding from GPL license")
	
	// Verify license agent finding
	foundLicenseFinding := false
	for _, result := range analysisResp.Results {
		if agentName, ok := result["agent_name"].(string); ok && agentName == "License Agent" {
			if finding, ok := result["finding"].(string); ok && strings.Contains(finding, "GPL-3.0-only") {
				foundLicenseFinding = true
				assert.Equal(t, "High", result["severity"], "GPL license should have High severity")
				t.Logf("✓ Found expected license finding: %s", finding)
				break
			}
		}
	}
	assert.True(t, foundLicenseFinding, "Should find GPL license compliance issue")

	t.Logf("✓ License analysis completed successfully")

	// Step 4: Analyze SBOM with multiple agents enabled
	t.Log("Step 4: Analyzing SBOM with multiple agents...")
	
	multiAgentURL := fmt.Sprintf("%s/api/v1/sboms/%s/analyze?enable-ai-health-check=true&enable-vuln-scan=true", 
		ts.Server.URL, sbomID)
	req, err = http.NewRequest("POST", multiAgentURL, nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Multi-agent analysis should return 200 OK")

	respBody, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var multiAgentResp AnalysisResponse
	err = json.Unmarshal(respBody, &multiAgentResp)
	require.NoError(t, err)

	// Verify multi-agent response
	assert.Equal(t, sbomID, multiAgentResp.SBOMID)
	
	// Check summary contains multiple agents
	if summary, ok := multiAgentResp.Summary["agents_run"].([]interface{}); ok {
		agentNames := make([]string, len(summary))
		for i, agent := range summary {
			agentNames[i] = agent.(string)
		}
		
		assert.Contains(t, agentNames, "License Agent", "Should include License Agent")
		assert.Contains(t, agentNames, "Dependency Health Agent", "Should include Dependency Health Agent")
		assert.Contains(t, agentNames, "Vulnerability Scanner", "Should include Vulnerability Scanner")
		
		t.Logf("✓ Multi-agent analysis ran agents: %v", agentNames)
	}

	t.Log("✓ Complete API workflow test passed successfully!")
}

func TestErrorHandling(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup()

	tests := []struct {
		name           string
		method         string
		url            string
		body           io.Reader
		contentType    string
		expectedStatus int
	}{
		{
			name:           "Invalid HTTP method for submit",
			method:         "GET",
			url:            "/api/v1/sboms",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Missing SBOM file",
			method:         "POST",
			url:            "/api/v1/sboms",
			body:           strings.NewReader(""),
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid SBOM ID for retrieval",
			method:         "GET",
			url:            "/api/v1/sboms/get?id=nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing SBOM ID for analysis",
			method:         "GET",
			url:            "/api/v1/sboms//analyze",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error
			
			if tt.body != nil {
				req, err = http.NewRequest(tt.method, ts.Server.URL+tt.url, tt.body)
			} else {
				req, err = http.NewRequest(tt.method, ts.Server.URL+tt.url, nil)
			}
			require.NoError(t, err)

			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, 
				"Expected status %d for %s", tt.expectedStatus, tt.name)
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Cleanup()

	const numConcurrentRequests = 5
	done := make(chan bool, numConcurrentRequests)
	errors := make(chan error, numConcurrentRequests)

	// Submit the same SBOM concurrently
	testSBOM := createTestSBOM()
	sbomJSON, err := json.Marshal(testSBOM)
	require.NoError(t, err)

	for i := 0; i < numConcurrentRequests; i++ {
		go func(requestID int) {
			defer func() { done <- true }()

			// Create multipart request
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)
			part, err := writer.CreateFormFile("sbom", fmt.Sprintf("test-sbom-%d.json", requestID))
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to create form file: %w", requestID, err)
				return
			}

			_, err = part.Write(sbomJSON)
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to write SBOM data: %w", requestID, err)
				return
			}

			err = writer.Close()
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to close writer: %w", requestID, err)
				return
			}

			// Submit SBOM
			req, err := http.NewRequest("POST", ts.Server.URL+"/api/v1/sboms", &requestBody)
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to create request: %w", requestID, err)
				return
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			client := &http.Client{Timeout: 30 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to send request: %w", requestID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				errors <- fmt.Errorf("request %d: expected status 201, got %d", requestID, resp.StatusCode)
				return
			}

			t.Logf("✓ Concurrent request %d completed successfully", requestID)
		}(i)
	}

	// Wait for all requests to complete
	completedRequests := 0
	for completedRequests < numConcurrentRequests {
		select {
		case <-done:
			completedRequests++
		case err := <-errors:
			t.Errorf("Concurrent request failed: %v", err)
			completedRequests++
		case <-time.After(60 * time.Second):
			t.Fatal("Concurrent requests timed out")
		}
	}

	t.Logf("✓ All %d concurrent requests completed", numConcurrentRequests)
}

// createTestSBOM creates a test SBOM with components that will trigger license findings
func createTestSBOM() map[string]interface{} {
	return map[string]interface{}{
		"bomFormat":    "CycloneDX",
		"specVersion":  "1.4",
		"serialNumber": "urn:uuid:integration-test-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		"version":      1,
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"component": map[string]interface{}{
				"type":    "application",
				"name":    "Test Application",
				"version": "1.0.0",
			},
		},
		"components": []map[string]interface{}{
			{
				"type":    "library",
				"name":    "express",
				"version": "4.18.2",
				"purl":    "pkg:npm/express@4.18.2",
				"licenses": []map[string]interface{}{
					{
						"license": map[string]interface{}{
							"id": "MIT",
						},
					},
				},
			},
			{
				"type":    "library",
				"name":    "copyleft-library",
				"version": "2.1.0",
				"purl":    "pkg:npm/copyleft-library@2.1.0",
				"licenses": []map[string]interface{}{
					{
						"license": map[string]interface{}{
							"id": "GPL-3.0-only",
						},
					},
				},
			},
			{
				"type":    "library",
				"name":    "another-gpl-lib",
				"version": "1.0.0",
				"purl":    "pkg:npm/another-gpl-lib@1.0.0",
				"licenses": []map[string]interface{}{
					{
						"license": map[string]interface{}{
							"id": "AGPL-3.0-only",
						},
					},
				},
			},
		},
	}
}