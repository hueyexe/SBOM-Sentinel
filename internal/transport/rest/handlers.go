// Package rest provides HTTP handlers for the SBOM Sentinel REST API.
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

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

// writeErrorResponse writes a standardized error response.
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Error:   errorType,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}