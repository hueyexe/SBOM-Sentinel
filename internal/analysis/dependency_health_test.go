package analysis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hueyexe/SBOM-Sentinel/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDependencyHealthAgent_Name(t *testing.T) {
	agent := NewDependencyHealthAgent()
	assert.Equal(t, "Dependency Health Agent", agent.Name())
}

func TestDependencyHealthAgent_Analyze(t *testing.T) {
	tests := []struct {
		name               string
		sbom               core.SBOM
		mockResponse       string
		mockStatusCode     int
		expectedCount      int
		expectedSeverities []string
		shouldReturnError  bool
	}{
		{
			name: "Component with risk identified",
			sbom: core.SBOM{
				ID:   "test-1",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "deprecated-library",
						Version: "1.0.0",
					},
				},
			},
			mockResponse:       `{"response": "This component is deprecated and no longer maintained."}`,
			mockStatusCode:     http.StatusOK,
			expectedCount:      1,
			expectedSeverities: []string{"Medium"},
			shouldReturnError:  false,
		},
		{
			name: "Component with no risk",
			sbom: core.SBOM{
				ID:   "test-2",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "healthy-library",
						Version: "2.1.0",
					},
				},
			},
			mockResponse:       `{"response": "This is a well-maintained and actively developed project."}`,
			mockStatusCode:     http.StatusOK,
			expectedCount:      0,
			expectedSeverities: []string{},
			shouldReturnError:  false,
		},
		{
			name: "Multiple components with mixed health",
			sbom: core.SBOM{
				ID:   "test-3",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "unmaintained-lib",
						Version: "0.9.0",
					},
					{
						Name:    "good-lib",
						Version: "3.0.0",
					},
				},
			},
			mockResponse:       `{"response": "unmaintained"}`,
			mockStatusCode:     http.StatusOK,
			expectedCount:      2, // Both will be checked, but only one will be flagged
			expectedSeverities: []string{"Medium"},
			shouldReturnError:  false,
		},
		{
			name: "Components without name or version",
			sbom: core.SBOM{
				ID:   "test-4",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "",
						Version: "1.0.0",
					},
					{
						Name:    "valid-component",
						Version: "",
					},
					{
						Name:    "complete-component",
						Version: "1.0.0",
					},
				},
			},
			mockResponse:       `{"response": "This is well maintained."}`,
			mockStatusCode:     http.StatusOK,
			expectedCount:      0,
			expectedSeverities: []string{},
			shouldReturnError:  false,
		},
		{
			name: "Empty SBOM",
			sbom: core.SBOM{
				ID:         "test-5",
				Name:       "Empty SBOM",
				Components: []core.Component{},
			},
			mockResponse:       "",
			mockStatusCode:     http.StatusOK,
			expectedCount:      0,
			expectedSeverities: []string{},
			shouldReturnError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			callCount := 0
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)

				// Alternate responses for multiple component tests
				if tt.name == "Multiple components with mixed health" {
					if callCount == 1 {
						w.Write([]byte(`{"response": "This component is unmaintained and deprecated."}`))
					} else {
						w.Write([]byte(`{"response": "This is a healthy, well-maintained project."}`))
					}
				} else {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer mockServer.Close()

			// Create agent with custom Ollama URL
			agent := NewDependencyHealthAgent()
			agent.ollamaURL = mockServer.URL

			ctx := context.Background()
			results, err := agent.Analyze(ctx, tt.sbom)

			if tt.shouldReturnError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Filter results based on expected count logic
			actualRiskyResults := []core.AnalysisResult{}
			for _, result := range results {
				if result.AgentName == "Dependency Health Agent" {
					actualRiskyResults = append(actualRiskyResults, result)
				}
			}

			// For multiple components test, we expect only the risky one
			if tt.name == "Multiple components with mixed health" {
				assert.Equal(t, 1, len(actualRiskyResults))
			} else {
				assert.Equal(t, tt.expectedCount, len(actualRiskyResults))
			}

			// Check result properties
			for i, result := range actualRiskyResults {
				assert.Equal(t, "Dependency Health Agent", result.AgentName)
				if i < len(tt.expectedSeverities) {
					assert.Equal(t, tt.expectedSeverities[i], result.Severity)
				}
				assert.NotEmpty(t, result.Finding)
			}
		})
	}
}

func TestDependencyHealthAgent_generatePrompt(t *testing.T) {
	agent := NewDependencyHealthAgent()
	component := core.Component{
		Name:    "test-library",
		Version: "1.2.3",
	}

	prompt := agent.generatePrompt(component)

	assert.Contains(t, prompt, "test-library")
	assert.Contains(t, prompt, "1.2.3")
	assert.Contains(t, prompt, "actively maintained")
	assert.Contains(t, prompt, "deprecated")
	assert.Contains(t, prompt, "risky")
}

func TestDependencyHealthAgent_queryOllama(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		expectedResult string
		shouldError    bool
	}{
		{
			name:           "Successful response",
			mockResponse:   `{"response": "This is a test response."}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "This is a test response.",
			shouldError:    false,
		},
		{
			name:           "Response with whitespace",
			mockResponse:   `{"response": "  \n  Test response with whitespace  \n  "}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "Test response with whitespace",
			shouldError:    false,
		},
		{
			name:           "Server error",
			mockResponse:   `{"error": "Internal server error"}`,
			mockStatusCode: http.StatusInternalServerError,
			expectedResult: "",
			shouldError:    true,
		},
		{
			name:           "Invalid JSON response",
			mockResponse:   `{invalid json}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and headers
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify request body contains expected fields
				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				assert.NoError(t, err)
				assert.Equal(t, "llama3", reqBody["model"])
				assert.Contains(t, reqBody["prompt"], "test prompt")
				assert.Equal(t, false, reqBody["stream"])

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer mockServer.Close()

			// Create agent with custom Ollama URL
			agent := NewDependencyHealthAgent()
			agent.ollamaURL = mockServer.URL

			ctx := context.Background()
			result, err := agent.queryOllama(ctx, "test prompt")

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestDependencyHealthAgent_indicatesRisk(t *testing.T) {
	agent := NewDependencyHealthAgent()

	tests := []struct {
		name     string
		response string
		expected bool
	}{
		{
			name:     "Unmaintained keyword",
			response: "This project is unmaintained and should be avoided.",
			expected: true,
		},
		{
			name:     "Deprecated keyword",
			response: "This library is deprecated in favor of newer alternatives.",
			expected: true,
		},
		{
			name:     "Multiple risk keywords",
			response: "This component is both outdated and has security issues.",
			expected: true,
		},
		{
			name:     "Case insensitive detection",
			response: "This project is UNMAINTAINED and RISKY to use.",
			expected: true,
		},
		{
			name:     "End of life detection",
			response: "This software has reached its end of life (EOL).",
			expected: true,
		},
		{
			name:     "Healthy project response",
			response: "This is a well-maintained, actively developed project with regular updates.",
			expected: false,
		},
		{
			name:     "Neutral response",
			response: "This is a standard library for web development.",
			expected: false,
		},
		{
			name:     "Empty response",
			response: "",
			expected: false,
		},
		{
			name:     "Archived project",
			response: "This repository has been archived by the owner.",
			expected: true,
		},
		{
			name:     "Security vulnerability mention",
			response: "This version has known security issues that were fixed in later versions.",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.indicatesRisk(tt.response)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDependencyHealthAgent_NetworkError(t *testing.T) {
	agent := NewDependencyHealthAgent()
	// Set an invalid URL to simulate network error
	agent.ollamaURL = "http://invalid-url:99999/api/generate"

	sbom := core.SBOM{
		ID:   "test",
		Name: "Test SBOM",
		Components: []core.Component{
			{
				Name:    "test-component",
				Version: "1.0.0",
			},
		},
	}

	ctx := context.Background()
	results, err := agent.Analyze(ctx, sbom)

	// Should not return error (graceful handling), but should have no results
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}
