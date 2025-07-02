package analysis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/vectordb"
	"github.com/stretchr/testify/assert"
)

func TestProactiveVulnerabilityAgent_Name(t *testing.T) {
	agent := NewProactiveVulnerabilityAgent()
	assert.Equal(t, "Proactive Vulnerability Agent", agent.Name())
}

func TestProactiveVulnerabilityAgent_Analyze(t *testing.T) {
	tests := []struct {
		name          string
		sbom          core.SBOM
		expectedCount int
	}{
		{
			name: "Component without name or version - should be skipped",
			sbom: core.SBOM{
				ID:   "test-1",
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
				},
			},
			expectedCount: 0,
		},
		{
			name: "Empty SBOM",
			sbom: core.SBOM{
				ID:         "test-2",
				Name:       "Empty SBOM",
				Components: []core.Component{},
			},
			expectedCount: 0,
		},
		{
			name: "Valid component - but will fail on network calls",
			sbom: core.SBOM{
				ID:   "test-3",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "test-component",
						Version: "1.0.0",
					},
				},
			},
			expectedCount: 0, // Will fail on embedding generation due to no mock server
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create agent
			agent := NewProactiveVulnerabilityAgent()

			// Mark as initialized to skip the initialization phase
			agent.initialized = true

			ctx := context.Background()
			results, err := agent.Analyze(ctx, tt.sbom)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(results))

			// Check each result
			for _, result := range results {
				assert.Equal(t, "Proactive Vulnerability Agent", result.AgentName)
				assert.Equal(t, "Medium", result.Severity)
				assert.NotEmpty(t, result.Finding)
			}
		})
	}
}

// Note: generateEmbedding is difficult to test in isolation due to hardcoded URLs
// It's tested implicitly through the integration tests in Analyze

func TestProactiveVulnerabilityAgent_queryLLM(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		mockResponse   string
		mockStatusCode int
		expectedResult string
		shouldError    bool
	}{
		{
			name:           "Successful LLM query",
			prompt:         "test security analysis prompt",
			mockResponse:   "This component has security vulnerabilities.",
			mockStatusCode: http.StatusOK,
			expectedResult: "This component has security vulnerabilities.",
			shouldError:    false,
		},
		{
			name:           "No concerns response - filtered out",
			prompt:         "test prompt",
			mockResponse:   "No relevant security concerns identified for this component.",
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			shouldError:    false,
		},
		{
			name:           "No vulnerabilities response - filtered out",
			prompt:         "test prompt",
			mockResponse:   "No vulnerabilities found in this component.",
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			shouldError:    false,
		},
		{
			name:           "Server error",
			prompt:         "test prompt",
			mockStatusCode: http.StatusInternalServerError,
			shouldError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				assert.NoError(t, err)
				assert.Equal(t, "llama3", reqBody["model"])
				assert.Equal(t, tt.prompt, reqBody["prompt"])
				assert.Equal(t, false, reqBody["stream"])

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)

				if tt.mockStatusCode == http.StatusOK {
					response := map[string]interface{}{
						"response": tt.mockResponse,
					}
					json.NewEncoder(w).Encode(response)
				}
			}))
			defer mockServer.Close()

			agent := NewProactiveVulnerabilityAgent()
			agent.ollamaURL = mockServer.URL

			ctx := context.Background()
			result, err := agent.queryLLM(ctx, tt.prompt)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestProactiveVulnerabilityAgent_analyzeWithLLM(t *testing.T) {
	component := core.Component{
		Name:    "test-component",
		Version: "1.0.0",
	}

	docs := []vectordb.Document{
		{
			ID:   "doc1",
			Text: "Security vulnerability in test-component version 1.0.0",
			Metadata: map[string]interface{}{
				"component": "test-component",
				"severity":  "High",
			},
		},
		{
			ID:   "doc2",
			Text: "Another security issue affecting test-component",
			Metadata: map[string]interface{}{
				"component": "test-component",
				"severity":  "Medium",
			},
		},
	}

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)
		
		prompt := reqBody["prompt"].(string)
		assert.Contains(t, prompt, "test-component")
		assert.Contains(t, prompt, "1.0.0")
		assert.Contains(t, prompt, "Security Intelligence Context")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		response := map[string]interface{}{
			"response": "Found potential security vulnerabilities in test-component.",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	agent := NewProactiveVulnerabilityAgent()
	agent.ollamaURL = mockServer.URL

	ctx := context.Background()
	result, err := agent.analyzeWithLLM(ctx, component, docs)

	assert.NoError(t, err)
	assert.Equal(t, "Found potential security vulnerabilities in test-component.", result)
}

func TestProactiveVulnerabilityAgent_NetworkError(t *testing.T) {
	agent := NewProactiveVulnerabilityAgent()
	// Set invalid URLs to simulate network errors
	agent.ollamaURL = "http://invalid-url:99999/api/generate"

	// Create a mock vector DB that will cause embedding errors
	agent.vectorDB = vectordb.NewMemoryVectorDB()
	// No mock documents, so search will return no results

	// Mark as initialized to skip initialization
	agent.initialized = true

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

// TestProactiveVulnerabilityAgent_InitializationError would be complex to test
// due to the dependency on external services and complex initialization flow
// The initialization logic is tested implicitly through integration tests