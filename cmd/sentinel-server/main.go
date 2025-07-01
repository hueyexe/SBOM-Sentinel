// Package main provides the entry point for the SBOM Sentinel server application.
// This binary serves as the REST API server for analyzing SBOM documents.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("SBOM Sentinel Server - Starting...")
	
	// TODO: Implement REST API handlers
	// TODO: Wire up dependency injection for ports and adapters
	// TODO: Add support for various SBOM format parsing via HTTP endpoints
	// TODO: Integrate analysis agents for API responses
	// TODO: Add proper configuration management
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"sbom-sentinel"}`))
	})
	
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}