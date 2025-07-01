// Package analysis provides dependency health analysis functionality for SBOM components.
package analysis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
)

// DependencyHealthAgent analyzes SBOM components for health and maintenance status using AI.
type DependencyHealthAgent struct {
	ollamaURL string
	model     string
	client    *http.Client
}

// OllamaRequest represents the request structure for Ollama API.
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response structure from Ollama API.
type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

// NewDependencyHealthAgent creates a new instance of DependencyHealthAgent.
func NewDependencyHealthAgent() *DependencyHealthAgent {
	return &DependencyHealthAgent{
		ollamaURL: "http://localhost:11434/api/generate",
		model:     "llama3",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the identifier for this analysis agent.
func (dha *DependencyHealthAgent) Name() string {
	return "Dependency Health Agent"
}

// Analyze examines the SBOM components for health and maintenance status using AI.
// It queries a local LLM via Ollama to assess each component's health.
func (dha *DependencyHealthAgent) Analyze(ctx context.Context, sbom core.SBOM) ([]core.AnalysisResult, error) {
	var results []core.AnalysisResult

	for _, component := range sbom.Components {
		// Skip components without name or version
		if component.Name == "" || component.Version == "" {
			continue
		}

		// Generate prompt for the LLM
		prompt := dha.generatePrompt(component)

		// Query the LLM
		response, err := dha.queryOllama(ctx, prompt)
		if err != nil {
			// Log error but continue with other components
			fmt.Printf("Warning: Failed to analyze component '%s': %v\n", component.Name, err)
			continue
		}

		// Check if the response indicates potential risk
		if dha.indicatesRisk(response) {
			result := core.AnalysisResult{
				AgentName: dha.Name(),
				Finding:   response,
				Severity:  "Medium",
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// generatePrompt creates a specific prompt for the LLM to assess component health.
func (dha *DependencyHealthAgent) generatePrompt(component core.Component) string {
	return fmt.Sprintf("Analyze the project health of the open-source component '%s' version '%s'. Based on public knowledge, is this project actively maintained, deprecated, or considered risky for other reasons? Answer in one sentence.",
		component.Name, component.Version)
}

// queryOllama sends a request to the Ollama API and returns the response.
func (dha *DependencyHealthAgent) queryOllama(ctx context.Context, prompt string) (string, error) {
	// Create request payload
	reqPayload := OllamaRequest{
		Model:  dha.model,
		Prompt: prompt,
		Stream: false,
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", dha.ollamaURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := dha.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}

// indicatesRisk checks if the LLM response indicates potential risk.
func (dha *DependencyHealthAgent) indicatesRisk(response string) bool {
	response = strings.ToLower(response)
	
	// Keywords that indicate potential risk
	riskKeywords := []string{
		"unmaintained",
		"deprecated",
		"risky",
		"outdated",
		"abandoned",
		"not maintained",
		"no longer maintained",
		"inactive",
		"archived",
		"obsolete",
		"discontinued",
		"end of life",
		"eol",
		"unsupported",
		"vulnerable",
		"security issues",
		"not recommended",
		"avoid",
		"stale",
		"dead project",
	}

	for _, keyword := range riskKeywords {
		if strings.Contains(response, keyword) {
			return true
		}
	}

	return false
}