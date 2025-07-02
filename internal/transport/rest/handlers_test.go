package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the storage.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Store(ctx context.Context, sbom core.SBOM) error {
	args := m.Called(ctx, sbom)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*core.SBOM, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.SBOM), args.Error(1)
}

func TestSubmitSBOMHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		setupRequest       func() (*http.Request, error)
		mockBehavior       func(*MockRepository)
		expectedStatusCode int
		expectedResponse   func(*testing.T, []byte)
	}{
		{
			name:   "Successful SBOM submission",
			method: "POST",
			setupRequest: func() (*http.Request, error) {
				// Create a valid CycloneDX SBOM
				sbomData := `{
					"bomFormat": "CycloneDX",
					"specVersion": "1.4",
					"serialNumber": "urn:uuid:test-12345",
					"version": 1,
					"metadata": {
						"timestamp": "2024-01-01T00:00:00Z"
					},
					"components": [
						{
							"type": "library",
							"name": "test-library",
							"version": "1.0.0",
							"purl": "pkg:npm/test-library@1.0.0",
							"licenses": [
								{
									"license": {
										"id": "MIT"
									}
								}
							]
						}
					]
				}`

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("sbom", "test.json")
				if err != nil {
					return nil, err
				}
				part.Write([]byte(sbomData))
				writer.Close()

				req := httptest.NewRequest("POST", "/api/v1/sboms", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("Store", mock.Anything, mock.AnythingOfType("core.SBOM")).Return(nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: func(t *testing.T, body []byte) {
				var response SubmitSBOMResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, "SBOM submitted successfully", response.Message)
			},
		},
		{
			name:   "Wrong HTTP method",
			method: "GET",
			setupRequest: func() (*http.Request, error) {
				return httptest.NewRequest("GET", "/api/v1/sboms", nil), nil
			},
			mockBehavior:       func(mockRepo *MockRepository) {},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "method_not_allowed", response.Error)
			},
		},
		{
			name:   "Missing SBOM file",
			method: "POST",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()

				req := httptest.NewRequest("POST", "/api/v1/sboms", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockBehavior:       func(mockRepo *MockRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "missing_file", response.Error)
			},
		},
		{
			name:   "Empty file",
			method: "POST",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("sbom", "empty.json")
				if err != nil {
					return nil, err
				}
				part.Write([]byte(""))
				writer.Close()

				req := httptest.NewRequest("POST", "/api/v1/sboms", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockBehavior:       func(mockRepo *MockRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "empty_file", response.Error)
			},
		},
		{
			name:   "Invalid SBOM format",
			method: "POST",
			setupRequest: func() (*http.Request, error) {
				invalidSbomData := `{"invalid": "json structure"}`

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("sbom", "invalid.json")
				if err != nil {
					return nil, err
				}
				part.Write([]byte(invalidSbomData))
				writer.Close()

				req := httptest.NewRequest("POST", "/api/v1/sboms", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockBehavior:       func(mockRepo *MockRepository) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "parse_error", response.Error)
			},
		},
		{
			name:   "Database storage error",
			method: "POST",
			setupRequest: func() (*http.Request, error) {
				sbomData := `{
					"bomFormat": "CycloneDX",
					"specVersion": "1.4",
					"serialNumber": "urn:uuid:test-12345",
					"version": 1,
					"components": []
				}`

				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("sbom", "test.json")
				if err != nil {
					return nil, err
				}
				part.Write([]byte(sbomData))
				writer.Close()

				req := httptest.NewRequest("POST", "/api/v1/sboms", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("Store", mock.Anything, mock.AnythingOfType("core.SBOM")).Return(errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "storage_error", response.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockRepository)
			tt.mockBehavior(mockRepo)

			// Create request
			req, err := tt.setupRequest()
			assert.NoError(t, err)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler and serve
			handler := SubmitSBOMHandler(mockRepo)
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			// Check response
			if tt.expectedResponse != nil {
				tt.expectedResponse(t, rr.Body.Bytes())
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetSBOMHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		queryParams        string
		mockBehavior       func(*MockRepository)
		expectedStatusCode int
		expectedResponse   func(*testing.T, []byte)
	}{
		{
			name:        "Successful SBOM retrieval",
			method:      "GET",
			queryParams: "?id=test-sbom-123",
			mockBehavior: func(mockRepo *MockRepository) {
				expectedSBOM := &core.SBOM{
					ID:   "test-sbom-123",
					Name: "Test SBOM",
					Components: []core.Component{
						{
							Name:    "test-component",
							Version: "1.0.0",
							License: "MIT",
						},
					},
				}
				mockRepo.On("FindByID", mock.Anything, "test-sbom-123").Return(expectedSBOM, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: func(t *testing.T, body []byte) {
				var sbom core.SBOM
				err := json.Unmarshal(body, &sbom)
				assert.NoError(t, err)
				assert.Equal(t, "test-sbom-123", sbom.ID)
				assert.Equal(t, "Test SBOM", sbom.Name)
				assert.Len(t, sbom.Components, 1)
			},
		},
		{
			name:   "Wrong HTTP method",
			method: "POST",
			mockBehavior: func(mockRepo *MockRepository) {
				// No expectations as method check happens first
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "method_not_allowed", response.Error)
			},
		},
		{
			name:        "Missing ID parameter",
			method:      "GET",
			queryParams: "",
			mockBehavior: func(mockRepo *MockRepository) {
				// No expectations as ID check happens first
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "missing_id", response.Error)
			},
		},
		{
			name:        "SBOM not found",
			method:      "GET",
			queryParams: "?id=nonexistent-sbom",
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("FindByID", mock.Anything, "nonexistent-sbom").Return(nil, nil)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "not_found", response.Error)
			},
		},
		{
			name:        "Database error",
			method:      "GET",
			queryParams: "?id=error-sbom",
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("FindByID", mock.Anything, "error-sbom").Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "storage_error", response.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockRepository)
			tt.mockBehavior(mockRepo)

			// Create request
			req := httptest.NewRequest(tt.method, "/api/v1/sboms"+tt.queryParams, nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler and serve
			handler := GetSBOMHandler(mockRepo)
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			// Check response
			if tt.expectedResponse != nil {
				tt.expectedResponse(t, rr.Body.Bytes())
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnalyzeSBOMHandler(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		urlPath            string
		queryParams        string
		mockBehavior       func(*MockRepository)
		expectedStatusCode int
		expectedResponse   func(*testing.T, []byte)
	}{
		{
			name:        "Successful analysis with license agent only",
			method:      "POST",
			urlPath:     "/api/v1/sboms/test-sbom-123/analyze",
			queryParams: "",
			mockBehavior: func(mockRepo *MockRepository) {
				testSBOM := &core.SBOM{
					ID:   "test-sbom-123",
					Name: "Test SBOM",
					Components: []core.Component{
						{
							Name:    "risky-component",
							Version: "1.0.0",
							License: "GPL-3.0-only",
						},
					},
				}
				mockRepo.On("FindByID", mock.Anything, "test-sbom-123").Return(testSBOM, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: func(t *testing.T, body []byte) {
				var response AnalysisResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "test-sbom-123", response.SBOMID)
				assert.Equal(t, 1, response.Summary.TotalFindings)
				assert.Contains(t, response.Summary.AgentsRun, "License Agent")
				assert.Len(t, response.Results, 1)
				assert.Equal(t, "License Agent", response.Results[0].AgentName)
			},
		},
		{
			name:        "Analysis with all agents enabled",
			method:      "POST",
			urlPath:     "/api/v1/sboms/test-sbom-456/analyze",
			queryParams: "?enable-ai-health-check=true&enable-proactive-scan=true&enable-vuln-scan=true",
			mockBehavior: func(mockRepo *MockRepository) {
				testSBOM := &core.SBOM{
					ID:   "test-sbom-456",
					Name: "Test SBOM",
					Components: []core.Component{
						{
							Name:    "safe-component",
							Version: "2.0.0",
							License: "MIT",
						},
					},
				}
				mockRepo.On("FindByID", mock.Anything, "test-sbom-456").Return(testSBOM, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: func(t *testing.T, body []byte) {
				var response AnalysisResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "test-sbom-456", response.SBOMID)
				// Should have all 4 agents in the summary
				assert.Len(t, response.Summary.AgentsRun, 4)
				assert.Contains(t, response.Summary.AgentsRun, "License Agent")
				assert.Contains(t, response.Summary.AgentsRun, "Dependency Health Agent")
				assert.Contains(t, response.Summary.AgentsRun, "Proactive Vulnerability Agent")
				assert.Contains(t, response.Summary.AgentsRun, "Vulnerability Scanner")
			},
		},
		{
			name:    "Wrong HTTP method",
			method:  "GET",
			urlPath: "/api/v1/sboms/test-sbom-123/analyze",
			mockBehavior: func(mockRepo *MockRepository) {
				// No expectations as method check happens first
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "method_not_allowed", response.Error)
			},
		},
		{
			name:    "Invalid URL path - missing SBOM ID",
			method:  "POST",
			urlPath: "/api/v1/sboms//analyze",  // Empty ID between slashes
			mockBehavior: func(mockRepo *MockRepository) {
				// No expectations as the empty ID is caught by validation first
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "missing_id", response.Error)
			},
		},
		{
			name:        "SBOM not found",
			method:      "POST",
			urlPath:     "/api/v1/sboms/nonexistent-sbom/analyze",
			queryParams: "",
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("FindByID", mock.Anything, "nonexistent-sbom").Return(nil, nil)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "not_found", response.Error)
			},
		},
		{
			name:        "Database error during SBOM retrieval",
			method:      "POST",
			urlPath:     "/api/v1/sboms/error-sbom/analyze",
			queryParams: "",
			mockBehavior: func(mockRepo *MockRepository) {
				mockRepo.On("FindByID", mock.Anything, "error-sbom").Return(nil, errors.New("database error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: func(t *testing.T, body []byte) {
				var response ErrorResponse
				err := json.Unmarshal(body, &response)
				assert.NoError(t, err)
				assert.Equal(t, "storage_error", response.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(MockRepository)
			tt.mockBehavior(mockRepo)

			// Create request
			fullURL := tt.urlPath
			if tt.queryParams != "" {
				fullURL += tt.queryParams
			}
			req := httptest.NewRequest(tt.method, fullURL, nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create handler and serve
			handler := AnalyzeSBOMHandler(mockRepo)
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			// Check response
			if tt.expectedResponse != nil {
				tt.expectedResponse(t, rr.Body.Bytes())
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateAnalysisSummary(t *testing.T) {
	tests := []struct {
		name             string
		results          []core.AnalysisResult
		agentsRun        []string
		expectedSummary  AnalysisSummary
	}{
		{
			name: "Multiple findings with different severities",
			results: []core.AnalysisResult{
				{AgentName: "License Agent", Finding: "GPL license found", Severity: "High"},
				{AgentName: "License Agent", Finding: "AGPL license found", Severity: "Critical"},
				{AgentName: "Vulnerability Scanner", Finding: "CVE found", Severity: "Medium"},
				{AgentName: "Dependency Health Agent", Finding: "Unmaintained library", Severity: "Medium"},
			},
			agentsRun: []string{"License Agent", "Vulnerability Scanner", "Dependency Health Agent"},
			expectedSummary: AnalysisSummary{
				TotalFindings: 4,
				FindingsBySeverity: map[string]int{
					"Critical": 1,
					"High":     1,
					"Medium":   2,
				},
				AgentsRun: []string{"License Agent", "Vulnerability Scanner", "Dependency Health Agent"},
			},
		},
		{
			name:    "No findings",
			results: []core.AnalysisResult{},
			agentsRun: []string{"License Agent"},
			expectedSummary: AnalysisSummary{
				TotalFindings:      0,
				FindingsBySeverity: map[string]int{},
				AgentsRun:          []string{"License Agent"},
			},
		},
		{
			name: "Single finding",
			results: []core.AnalysisResult{
				{AgentName: "License Agent", Finding: "MIT license found", Severity: "Low"},
			},
			agentsRun: []string{"License Agent"},
			expectedSummary: AnalysisSummary{
				TotalFindings: 1,
				FindingsBySeverity: map[string]int{
					"Low": 1,
				},
				AgentsRun: []string{"License Agent"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := generateAnalysisSummary(tt.results, tt.agentsRun)
			
			assert.Equal(t, tt.expectedSummary.TotalFindings, summary.TotalFindings)
			assert.Equal(t, tt.expectedSummary.FindingsBySeverity, summary.FindingsBySeverity)
			assert.Equal(t, tt.expectedSummary.AgentsRun, summary.AgentsRun)
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errorType      string
		message        string
		expectedStatus int
		expectedBody   ErrorResponse
	}{
		{
			name:           "Bad request error",
			statusCode:     http.StatusBadRequest,
			errorType:      "invalid_input",
			message:        "Input validation failed",
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Error:   "invalid_input",
				Message: "Input validation failed",
			},
		},
		{
			name:           "Internal server error",
			statusCode:     http.StatusInternalServerError,
			errorType:      "database_error",
			message:        "Database connection failed",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: ErrorResponse{
				Error:   "database_error",
				Message: "Database connection failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create response recorder
			rr := httptest.NewRecorder()

			// Call writeErrorResponse
			writeErrorResponse(rr, tt.statusCode, tt.errorType, tt.message)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body
			var response ErrorResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

// Helper function to create a valid multipart request for testing
func createMultipartRequest(fieldName, fileName, content string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	
	_, err = io.WriteString(part, content)
	if err != nil {
		return nil, err
	}
	
	writer.Close()
	
	req := httptest.NewRequest("POST", "/api/v1/sboms", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	return req, nil
}