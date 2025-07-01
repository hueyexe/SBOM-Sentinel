// Package main provides the entry point for the SBOM Sentinel server application.
// This binary serves as the REST API server for analyzing SBOM documents.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/database"
	"github.com/chrisclapham/SBOM-Sentinel/internal/transport/rest"
)

func main() {
	fmt.Println("SBOM Sentinel Server - Starting...")
	
	// Initialize SQLite database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./sentinel.db"
	}
	
	repo, err := database.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()
	
	fmt.Printf("Database initialized: %s\n", dbPath)
	
	// Configure routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"sbom-sentinel"}`))
	})
	
	// API v1 routes
	http.HandleFunc("/api/v1/sboms", rest.SubmitSBOMHandler(repo))
	http.HandleFunc("/api/v1/sboms/get", rest.GetSBOMHandler(repo))
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/v1/sboms       - Submit SBOM file")
	fmt.Println("  GET  /api/v1/sboms/get   - Retrieve SBOM by ID")
	fmt.Println("  GET  /health             - Health check")
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}