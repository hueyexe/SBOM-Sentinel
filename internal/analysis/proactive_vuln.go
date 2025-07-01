// Package analysis provides proactive vulnerability discovery using RAG pipeline.
package analysis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/vectordb"
)

// ProactiveVulnerabilityAgent analyzes SBOM components for potential vulnerabilities using RAG.
type ProactiveVulnerabilityAgent struct {
	vectorDB    *vectordb.MemoryVectorDB
	harvester   *vectordb.Harvester
	ollamaURL   string
	client      *http.Client
	initialized bool
}

// NewProactiveVulnerabilityAgent creates a new instance of ProactiveVulnerabilityAgent.
func NewProactiveVulnerabilityAgent() *ProactiveVulnerabilityAgent {
	vectorDB := vectordb.NewMemoryVectorDB()
	harvester := vectordb.NewHarvester(vectorDB)
	
	return &ProactiveVulnerabilityAgent{
		vectorDB:  vectorDB,
		harvester: harvester,
		ollamaURL: "http://localhost:11434/api/generate",
		client: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for RAG queries
		},
		initialized: false,
	}
}

// Name returns the identifier for this analysis agent.
func (pva *ProactiveVulnerabilityAgent) Name() string {
	return "Proactive Vulnerability Agent"
}

// Analyze examines the SBOM components for potential vulnerabilities using RAG pipeline.
func (pva *ProactiveVulnerabilityAgent) Analyze(ctx context.Context, sbom core.SBOM) ([]core.AnalysisResult, error) {
	// Initialize the vector database with mock security data if not already done
	if !pva.initialized {
		if err := pva.initializeSecurityIntelligence(ctx); err != nil {
			return nil, fmt.Errorf("failed to initialize security intelligence: %w", err)
		}
		pva.initialized = true
	}

	var results []core.AnalysisResult

	for _, component := range sbom.Components {
		// Skip components without name or version
		if component.Name == "" || component.Version == "" {
			continue
		}

		// Create embedding for the component query
		componentQuery := fmt.Sprintf("component %s version %s vulnerability security issue", component.Name, component.Version)
		queryEmbedding, err := pva.generateEmbedding(ctx, componentQuery)
		if err != nil {
			fmt.Printf("Warning: Failed to generate embedding for component '%s': %v\n", component.Name, err)
			continue
		}

		// Search for relevant security documents
		searchResults, err := pva.vectorDB.Search(queryEmbedding, 3) // Top 3 most relevant
		if err != nil {
			fmt.Printf("Warning: Failed to search vector DB for component '%s': %v\n", component.Name, err)
			continue
		}

		// Filter for relevant results with sufficient similarity
		var relevantDocs []vectordb.Document
		for _, result := range searchResults {
			if result.Similarity > 0.3 { // Only consider documents with >30% similarity
				relevantDocs = append(relevantDocs, result.Document)
			}
		}

		// If relevant documents found, query LLM for analysis
		if len(relevantDocs) > 0 {
			finding, err := pva.analyzeWithLLM(ctx, component, relevantDocs)
			if err != nil {
				fmt.Printf("Warning: Failed LLM analysis for component '%s': %v\n", component.Name, err)
				continue
			}

			if finding != "" {
				result := core.AnalysisResult{
					AgentName: pva.Name(),
					Finding:   finding,
					Severity:  "Medium", // RAG-discovered vulnerabilities are typically medium severity
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
}

// initializeSecurityIntelligence populates the vector database with security intelligence data.
func (pva *ProactiveVulnerabilityAgent) initializeSecurityIntelligence(ctx context.Context) error {
	fmt.Println("🔍 Initializing security intelligence database...")
	
	if err := pva.harvester.HarvestMockData(ctx); err != nil {
		return fmt.Errorf("failed to harvest security data: %w", err)
	}
	
	fmt.Printf("✅ Security intelligence database initialized with %d documents\n", pva.vectorDB.Size())
	return nil
}

// analyzeWithLLM uses the LLM to analyze component against relevant security documents.
func (pva *ProactiveVulnerabilityAgent) analyzeWithLLM(ctx context.Context, component core.Component, docs []vectordb.Document) (string, error) {
	// Build context from relevant documents
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Security Intelligence Context:\n")
	
	for i, doc := range docs {
		contextBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc.Text))
	}
	
	// Create prompt for LLM
	prompt := fmt.Sprintf(`Based on the security intelligence context provided, analyze if the component '%s' version '%s' has any potential security vulnerabilities or risks.

%s

Component to analyze: %s (version %s)

Instructions:
1. Look for any mentions of this specific component or similar components
2. Consider version compatibility and potential security issues
3. If you find relevant security concerns, summarize them in one sentence
4. If no relevant security issues are found, respond with "No relevant security concerns identified"

Response:`, component.Name, component.Version, contextBuilder.String(), component.Name, component.Version)

	return pva.queryLLM(ctx, prompt)
}



// queryLLM sends a query to the LLM and returns the response.
func (pva *ProactiveVulnerabilityAgent) queryLLM(ctx context.Context, prompt string) (string, error) {
	reqPayload := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", pva.ollamaURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pva.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	response := strings.TrimSpace(ollamaResp.Response)
	
	// Filter out "no concerns" responses
	if strings.Contains(strings.ToLower(response), "no relevant security concerns") ||
		strings.Contains(strings.ToLower(response), "no security issues") ||
		strings.Contains(strings.ToLower(response), "no vulnerabilities") {
		return "", nil
	}

	return response, nil
}

// generateEmbedding generates an embedding for the given text using Ollama.
func (pva *ProactiveVulnerabilityAgent) generateEmbedding(ctx context.Context, text string) ([]float64, error) {
	reqPayload := OllamaEmbeddingRequest{
		Model:  "llama3",
		Prompt: text,
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/embeddings", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pva.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Embedding, nil
}