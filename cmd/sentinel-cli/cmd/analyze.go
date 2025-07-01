// Package cmd provides the analyze command for parsing SBOM files.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/chrisclapham/SBOM-Sentinel/internal/ingestion"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [SBOM_FILE]",
	Short: "Analyze an SBOM file",
	Long: `Analyze a Software Bill of Materials (SBOM) file to extract component information.

Currently supports:
- CycloneDX JSON format

The command will parse the SBOM file and display information about the
components found within it.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	
	// Add flags specific to the analyze command
	analyzeCmd.Flags().StringP("format", "f", "auto", "SBOM format (auto, cyclonedx)")
	analyzeCmd.Flags().BoolP("summary", "s", false, "Show only summary information")
}

// runAnalyze executes the analyze command
func runAnalyze(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	
	// Check if verbose flag is set
	verbose, _ := cmd.Flags().GetBool("verbose")
	summary, _ := cmd.Flags().GetBool("summary")
	format, _ := cmd.Flags().GetString("format")
	
	if verbose {
		fmt.Printf("Analyzing SBOM file: %s\n", filePath)
		fmt.Printf("Format: %s\n", format)
	}
	
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()
	
	// For now, we only support CycloneDX JSON format
	// In the future, we could auto-detect format or support multiple parsers
	parser := ingestion.NewCycloneDXParser()
	
	// Parse the SBOM
	sbom, err := parser.Parse(file)
	if err != nil {
		return fmt.Errorf("failed to parse SBOM: %w", err)
	}
	
	// Display results
	fmt.Printf("✅ Successfully parsed SBOM: %s\n", sbom.Name)
	fmt.Printf("📦 Found %d components\n", len(sbom.Components))
	
	if !summary {
		fmt.Printf("\n📋 SBOM Details:\n")
		fmt.Printf("   ID: %s\n", sbom.ID)
		fmt.Printf("   Name: %s\n", sbom.Name)
		
		if len(sbom.Metadata) > 0 {
			fmt.Printf("\n🏷️  Metadata:\n")
			for key, value := range sbom.Metadata {
				fmt.Printf("   %s: %s\n", key, value)
			}
		}
		
		if len(sbom.Components) > 0 {
			fmt.Printf("\n🔍 Components:\n")
			for i, component := range sbom.Components {
				if i >= 10 && !verbose {
					fmt.Printf("   ... and %d more components (use --verbose to see all)\n", len(sbom.Components)-10)
					break
				}
				
				fmt.Printf("   • %s", component.Name)
				if component.Version != "" {
					fmt.Printf(" v%s", component.Version)
				}
				if component.License != "" {
					fmt.Printf(" (%s)", component.License)
				}
				fmt.Printf("\n")
				
				if verbose && component.PURL != "" {
					fmt.Printf("     PURL: %s\n", component.PURL)
				}
			}
		}
	}
	
	return nil
}