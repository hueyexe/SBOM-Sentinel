// Package cmd provides the analyze command for parsing SBOM files.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hueyexe/SBOM-Sentinel/internal/analysis"
	"github.com/hueyexe/SBOM-Sentinel/internal/core"
	"github.com/hueyexe/SBOM-Sentinel/internal/ingestion"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [SBOM_FILE]",
	Short: "Analyze an SBOM file",
	Long: `Analyze a Software Bill of Materials (SBOM) file to extract component information.

Currently supports:
- CycloneDX JSON format
- License compliance analysis
- AI-powered dependency health analysis (with --enable-ai-health-check)
- Proactive vulnerability discovery using RAG (with --enable-proactive-scan)

The command will parse the SBOM file and display information about the
components found within it, along with any security or compliance findings.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Add flags specific to the analyze command
	analyzeCmd.Flags().StringP("format", "f", "auto", "SBOM format (auto, cyclonedx)")
	analyzeCmd.Flags().BoolP("summary", "s", false, "Show only summary information")
	analyzeCmd.Flags().Bool("enable-ai-health-check", false, "Enable AI-powered dependency health analysis (requires Ollama)")
	analyzeCmd.Flags().Bool("enable-proactive-scan", false, "Enable proactive vulnerability discovery using RAG (requires Ollama)")
	analyzeCmd.Flags().Bool("enable-vuln-scan", false, "Enable known vulnerability scanning using OSV.dev database")
}

// runAnalyze executes the analyze command
func runAnalyze(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Check if verbose flag is set
	verbose, _ := cmd.Flags().GetBool("verbose")
	summary, _ := cmd.Flags().GetBool("summary")
	format, _ := cmd.Flags().GetString("format")
	enableAIHealthCheck, _ := cmd.Flags().GetBool("enable-ai-health-check")
	enableProactiveScan, _ := cmd.Flags().GetBool("enable-proactive-scan")
	enableVulnScan, _ := cmd.Flags().GetBool("enable-vuln-scan")

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

	// Run analysis agents
	ctx := context.Background()
	var allAnalysisResults []core.AnalysisResult

	// Run license analysis
	licenseAgent := analysis.NewLicenseAgent()

	if verbose {
		fmt.Printf("🔍 Running license analysis...\n")
	}

	licenseResults, err := licenseAgent.Analyze(ctx, *sbom)
	if err != nil {
		return fmt.Errorf("failed to run license analysis: %w", err)
	}
	allAnalysisResults = append(allAnalysisResults, licenseResults...)

	// Run AI health check if enabled
	if enableAIHealthCheck {
		healthAgent := analysis.NewDependencyHealthAgent()

		if verbose {
			fmt.Printf("🤖 Running AI-powered dependency health analysis...\n")
		}

		healthResults, err := healthAgent.Analyze(ctx, *sbom)
		if err != nil {
			fmt.Printf("Warning: AI health analysis failed: %v\n", err)
		} else {
			allAnalysisResults = append(allAnalysisResults, healthResults...)
		}
	}

	// Run proactive vulnerability scan if enabled
	if enableProactiveScan {
		proactiveAgent := analysis.NewProactiveVulnerabilityAgent()

		if verbose {
			fmt.Printf("🔍 Running proactive vulnerability discovery using RAG...\n")
		}

		proactiveResults, err := proactiveAgent.Analyze(ctx, *sbom)
		if err != nil {
			fmt.Printf("Warning: Proactive vulnerability scan failed: %v\n", err)
		} else {
			allAnalysisResults = append(allAnalysisResults, proactiveResults...)
		}
	}

	// Run vulnerability scan if enabled
	if enableVulnScan {
		vulnAgent := analysis.NewVulnerabilityScanningAgent()

		if verbose {
			fmt.Printf("🔍 Running known vulnerability scan using OSV.dev...\n")
		}

		vulnResults, err := vulnAgent.Analyze(ctx, *sbom)
		if err != nil {
			fmt.Printf("Warning: Vulnerability scan failed: %v\n", err)
		} else {
			allAnalysisResults = append(allAnalysisResults, vulnResults...)
		}
	}

	// Display analysis results if any findings were detected
	if len(allAnalysisResults) > 0 {
		fmt.Printf("\n🔬 Analysis Results:\n")
		fmt.Printf("   Found %d issues:\n\n", len(allAnalysisResults))

		for i, result := range allAnalysisResults {
			severityIcon := getSeverityIcon(result.Severity)
			fmt.Printf("   %d. %s [%s] %s\n", i+1, severityIcon, result.Severity, result.AgentName)
			fmt.Printf("      %s\n", result.Finding)
			if i < len(allAnalysisResults)-1 {
				fmt.Printf("\n")
			}
		}
	} else {
		fmt.Printf("\n✅ Analysis Complete: No issues detected\n")
		if !enableAIHealthCheck {
			fmt.Printf("   💡 Tip: Use --enable-ai-health-check for AI-powered dependency health analysis\n")
		}
		if !enableProactiveScan {
			fmt.Printf("   🔍 Tip: Use --enable-proactive-scan for proactive vulnerability discovery using RAG\n")
		}
		if !enableVulnScan {
			fmt.Printf("   🛡️  Tip: Use --enable-vuln-scan for known vulnerability scanning using OSV.dev\n")
		}
	}

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

// getSeverityIcon returns an appropriate emoji icon for the given severity level.
func getSeverityIcon(severity string) string {
	switch severity {
	case "Critical":
		return "🚨"
	case "High":
		return "🔴"
	case "Medium":
		return "🟡"
	case "Low":
		return "🟢"
	default:
		return "⚠️"
	}
}
