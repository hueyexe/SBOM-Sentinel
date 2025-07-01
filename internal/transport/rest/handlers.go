// Package rest provides HTTP handlers for the SBOM Sentinel REST API.
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chrisclapham/SBOM-Sentinel/internal/analysis"
	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
	"github.com/chrisclapham/SBOM-Sentinel/internal/ingestion"
	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/storage"
)

// SubmitSBOMResponse represents the JSON response for SBOM submission.
type SubmitSBOMResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// AnalysisResponse represents the JSON response for SBOM analysis.
type AnalysisResponse struct {
	SBOMID  string                `json:"sbom_id"`
	Results []core.AnalysisResult `json:"results"`
	Summary AnalysisSummary       `json:"summary"`
}

// AnalysisSummary provides a summary of the analysis results.
type AnalysisSummary struct {
	TotalFindings   int            `json:"total_findings"`
	FindingsBySeverity map[string]int `json:"findings_by_severity"`
	AgentsRun       []string       `json:"agents_run"`
}

// SubmitSBOMHandler creates an HTTP handler for submitting SBOM files.
// It expects a multipart/form-data request with an SBOM file.
func SubmitSBOMHandler(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST method is allowed")
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Parse multipart form (32MB max memory)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "invalid_form", "Failed to parse multipart form")
			return
		}

		// Get the uploaded file
		file, header, err := r.FormFile("sbom")
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "missing_file", "SBOM file is required. Please upload a file with the 'sbom' field name")
			return
		}
		defer file.Close()

		// Validate file type (optional - could check file extension)
		if header.Size == 0 {
			writeErrorResponse(w, http.StatusBadRequest, "empty_file", "Uploaded file is empty")
			return
		}

		// Create parser instance
		parser := ingestion.NewCycloneDXParser()

		// Parse the SBOM file
		sbom, err := parser.Parse(file)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "parse_error", fmt.Sprintf("Failed to parse SBOM file: %v", err))
			return
		}

		// Store the SBOM in the database
		ctx := r.Context()
		err = repo.Store(ctx, *sbom)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "storage_error", fmt.Sprintf("Failed to store SBOM: %v", err))
			return
		}

		// Return success response
		response := SubmitSBOMResponse{
			ID:      sbom.ID,
			Message: "SBOM submitted successfully",
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetSBOMHandler creates an HTTP handler for retrieving SBOM by ID.
func GetSBOMHandler(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET method is allowed")
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Extract ID from URL path
		// For simplicity, we'll expect the ID as a query parameter
		// In a real application, you might use a router like gorilla/mux
		id := r.URL.Query().Get("id")
		if id == "" {
			writeErrorResponse(w, http.StatusBadRequest, "missing_id", "SBOM ID is required as query parameter")
			return
		}

		// Retrieve SBOM from database
		ctx := r.Context()
		sbom, err := repo.FindByID(ctx, id)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "storage_error", fmt.Sprintf("Failed to retrieve SBOM: %v", err))
			return
		}

		if sbom == nil {
			writeErrorResponse(w, http.StatusNotFound, "not_found", "SBOM not found")
			return
		}

		// Return the SBOM
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sbom)
	}
}

// AnalyzeSBOMHandler creates an HTTP handler for analyzing stored SBOMs.
// It expects a POST request to /api/v1/sboms/{id}/analyze with optional query parameters.
func AnalyzeSBOMHandler(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only POST method is allowed")
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Extract SBOM ID from URL path
		// Expected format: /api/v1/sboms/{id}/analyze
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 4 || pathParts[3] == "" {
			writeErrorResponse(w, http.StatusBadRequest, "missing_id", "SBOM ID is required in URL path")
			return
		}
		sbomID := pathParts[3]

		// Check for AI health check flag
		enableAIHealthCheck := r.URL.Query().Get("enable-ai-health-check") == "true"
		// Check for proactive scan flag
		enableProactiveScan := r.URL.Query().Get("enable-proactive-scan") == "true"

		// Retrieve SBOM from database
		ctx := r.Context()
		sbom, err := repo.FindByID(ctx, sbomID)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "storage_error", fmt.Sprintf("Failed to retrieve SBOM: %v", err))
			return
		}

		if sbom == nil {
			writeErrorResponse(w, http.StatusNotFound, "not_found", "SBOM not found")
			return
		}

		// Run analysis agents
		var allResults []core.AnalysisResult
		var agentsRun []string

		// Run license analysis
		licenseAgent := analysis.NewLicenseAgent()
		licenseResults, err := licenseAgent.Analyze(ctx, *sbom)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "analysis_error", fmt.Sprintf("License analysis failed: %v", err))
			return
		}
		allResults = append(allResults, licenseResults...)
		agentsRun = append(agentsRun, licenseAgent.Name())

		// Run AI health check if enabled
		if enableAIHealthCheck {
			healthAgent := analysis.NewDependencyHealthAgent()
			healthResults, err := healthAgent.Analyze(ctx, *sbom)
			if err != nil {
				// Log warning but don't fail the entire analysis
				fmt.Printf("Warning: AI health analysis failed: %v\n", err)
			} else {
				allResults = append(allResults, healthResults...)
			}
			agentsRun = append(agentsRun, healthAgent.Name())
		}

		// Run proactive vulnerability scan if enabled
		if enableProactiveScan {
			proactiveAgent := analysis.NewProactiveVulnerabilityAgent()
			proactiveResults, err := proactiveAgent.Analyze(ctx, *sbom)
			if err != nil {
				// Log warning but don't fail the entire analysis
				fmt.Printf("Warning: Proactive vulnerability scan failed: %v\n", err)
			} else {
				allResults = append(allResults, proactiveResults...)
			}
			agentsRun = append(agentsRun, proactiveAgent.Name())
		}

		// Generate summary
		summary := generateAnalysisSummary(allResults, agentsRun)

		// Create response
		response := AnalysisResponse{
			SBOMID:  sbomID,
			Results: allResults,
			Summary: summary,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// generateAnalysisSummary creates a summary of analysis results.
func generateAnalysisSummary(results []core.AnalysisResult, agentsRun []string) AnalysisSummary {
	findingsBySeverity := make(map[string]int)
	
	for _, result := range results {
		findingsBySeverity[result.Severity]++
	}

	return AnalysisSummary{
		TotalFindings:      len(results),
		FindingsBySeverity: findingsBySeverity,
		AgentsRun:          agentsRun,
	}
}

// writeErrorResponse writes a standardized error response.
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Error:   errorType,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}