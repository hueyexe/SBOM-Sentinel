// Package main provides the entry point for the SBOM Sentinel CLI application.
// This binary serves as the command-line interface for analyzing SBOM documents.
package main

import (
	"fmt"
	"os"

	"github.com/hueyexe/SBOM-Sentinel/cmd/sentinel-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
