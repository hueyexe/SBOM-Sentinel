// Package vectordb provides harvesting functionality for security intelligence data.
package vectordb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SecurityIntelligence represents mock security intelligence data.
type SecurityIntelligence struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Component   string `json:"component"`
	Version     string `json:"version"`
	Severity    string `json:"severity"`
	Source      string `json:"source"`
	Date        string `json:"date"`
}

// Harvester handles the collection and processing of security intelligence data.
type Harvester struct {
	vectorDB    *MemoryVectorDB
	ollamaURL   string
	client      *http.Client
}

// NewHarvester creates a new Harvester instance.
func NewHarvester(vectorDB *MemoryVectorDB) *Harvester {
	return &Harvester{
		vectorDB:  vectorDB,
		ollamaURL: "http://localhost:11434/api/embeddings",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HarvestMockData creates and processes mock security intelligence data.
func (h *Harvester) HarvestMockData(ctx context.Context) error {
	mockData := h.generateMockSecurityData()
	
	for _, intelligence := range mockData {
		// Create document text from intelligence data
		docText := fmt.Sprintf("Title: %s. Description: %s. Component: %s, Version: %s. Severity: %s. Source: %s.",
			intelligence.Title,
			intelligence.Description,
			intelligence.Component,
			intelligence.Version,
			intelligence.Severity,
			intelligence.Source)
		
		// Generate embedding for the document
		embedding, err := h.generateEmbedding(ctx, docText)
		if err != nil {
			fmt.Printf("Warning: Failed to generate embedding for document %s: %v\n", intelligence.ID, err)
			continue
		}
		
		// Create document and add to vector database
		doc := Document{
			ID:     intelligence.ID,
			Text:   docText,
			Vector: embedding,
			Metadata: map[string]interface{}{
				"component": intelligence.Component,
				"version":   intelligence.Version,
				"severity":  intelligence.Severity,
				"source":    intelligence.Source,
				"date":      intelligence.Date,
				"title":     intelligence.Title,
			},
		}
		
		if err := h.vectorDB.Add(doc); err != nil {
			fmt.Printf("Warning: Failed to add document to vector DB: %v\n", err)
		}
	}
	
	fmt.Printf("Successfully harvested %d security intelligence documents\n", len(mockData))
	return nil
}

// generateMockSecurityData creates mock security intelligence data.
func (h *Harvester) generateMockSecurityData() []SecurityIntelligence {
	return []SecurityIntelligence{
		{
			ID:          "vuln-001",
			Title:       "Deserialization Vulnerability in acme-serializer",
			Description: "A new deserialization issue is being discussed for the 'acme-serializer' library version 1.2.3, allowing potential remote code execution. Researchers have identified unsafe deserialization patterns that could be exploited.",
			Component:   "acme-serializer",
			Version:     "1.2.3",
			Severity:    "Critical",
			Source:      "Security Mailing List",
			Date:        "2024-01-15",
		},
		{
			ID:          "vuln-002",
			Title:       "Memory Leak in data-processor",
			Description: "Security researchers are reporting memory leak issues in data-processor version 2.1.0 that could lead to denial of service attacks. The leak occurs during heavy processing workloads.",
			Component:   "data-processor",
			Version:     "2.1.0",
			Severity:    "High",
			Source:      "Research Blog",
			Date:        "2024-01-14",
		},
		{
			ID:          "vuln-003",
			Title:       "SQL Injection in database-connector",
			Description: "A potential SQL injection vulnerability has been identified in database-connector library version 3.4.1. The issue affects parameterized query handling in certain edge cases.",
			Component:   "database-connector",
			Version:     "3.4.1",
			Severity:    "High",
			Source:      "Security Forum",
			Date:        "2024-01-13",
		},
		{
			ID:          "vuln-004",
			Title:       "Path Traversal in file-manager",
			Description: "Discussions on security forums indicate a path traversal vulnerability in file-manager version 1.8.0 that allows access to files outside the intended directory structure.",
			Component:   "file-manager",
			Version:     "1.8.0",
			Severity:    "Medium",
			Source:      "Security Forum",
			Date:        "2024-01-12",
		},
		{
			ID:          "vuln-005",
			Title:       "XSS Vulnerability in web-utils",
			Description: "Cross-site scripting vulnerability discovered in web-utils version 2.3.4. The issue affects input sanitization functions and could allow malicious script execution.",
			Component:   "web-utils",
			Version:     "2.3.4",
			Severity:    "Medium",
			Source:      "Security Blog",
			Date:        "2024-01-11",
		},
		{
			ID:          "vuln-006",
			Title:       "Privilege Escalation in auth-service",
			Description: "Research indicates a privilege escalation issue in auth-service library version 4.2.1 where normal users can gain administrative privileges through token manipulation.",
			Component:   "auth-service",
			Version:     "4.2.1",
			Severity:    "Critical",
			Source:      "Research Paper",
			Date:        "2024-01-10",
		},
		{
			ID:          "vuln-007",
			Title:       "Buffer Overflow in image-processor",
			Description: "Security mailing lists are discussing a buffer overflow vulnerability in image-processor version 1.5.2 when processing specially crafted image files.",
			Component:   "image-processor",
			Version:     "1.5.2",
			Severity:    "High",
			Source:      "Security Mailing List",
			Date:        "2024-01-09",
		},
		{
			ID:          "vuln-008",
			Title:       "Information Disclosure in logger-util",
			Description: "Researchers have identified an information disclosure vulnerability in logger-util version 0.9.1 that may leak sensitive data in log files under certain configurations.",
			Component:   "logger-util",
			Version:     "0.9.1",
			Severity:    "Low",
			Source:      "Research Blog",
			Date:        "2024-01-08",
		},
	}
}

// OllamaEmbeddingRequest represents the request structure for Ollama embeddings API.
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the response structure from Ollama embeddings API.
type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// generateEmbedding generates an embedding for the given text using Ollama.
func (h *Harvester) generateEmbedding(ctx context.Context, text string) ([]float64, error) {
	reqPayload := OllamaEmbeddingRequest{
		Model:  "llama3",
		Prompt: text,
	}
	
	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", h.ollamaURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := h.client.Do(req)
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