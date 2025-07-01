// Package cmd provides the CLI commands for SBOM Sentinel.
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sentinel-cli",
	Short: "SBOM Sentinel - Advanced Software Bill of Materials analysis",
	Long: `SBOM Sentinel is an advanced Software Bill of Materials (SBOM) analysis platform.
It provides deep, contextual intelligence on software supply chain security,
license compliance, and dependency health using local AI engines.

This CLI tool allows you to analyze SBOM documents in various formats
including CycloneDX and SPDX.`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags here if needed
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}