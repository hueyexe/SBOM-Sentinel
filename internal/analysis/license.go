// Package analysis provides license analysis functionality for SBOM components.
package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
)

// LicenseAgent analyzes SBOM components for high-risk copyleft licenses.
type LicenseAgent struct {
	highRiskLicenses map[string]string
}

// NewLicenseAgent creates a new instance of LicenseAgent with predefined high-risk licenses.
func NewLicenseAgent() *LicenseAgent {
	// Define high-risk copyleft licenses that may pose compliance risks
	highRiskLicenses := map[string]string{
		"AGPL-3.0-only":        "GNU Affero General Public License v3.0 only",
		"AGPL-3.0-or-later":    "GNU Affero General Public License v3.0 or later",
		"GPL-2.0-only":         "GNU General Public License v2.0 only",
		"GPL-2.0-or-later":     "GNU General Public License v2.0 or later",
		"GPL-3.0-only":         "GNU General Public License v3.0 only",
		"GPL-3.0-or-later":     "GNU General Public License v3.0 or later",
		"LGPL-2.1-only":        "GNU Lesser General Public License v2.1 only",
		"LGPL-2.1-or-later":    "GNU Lesser General Public License v2.1 or later",
		"LGPL-3.0-only":        "GNU Lesser General Public License v3.0 only",
		"LGPL-3.0-or-later":    "GNU Lesser General Public License v3.0 or later",
		"EUPL-1.1":             "European Union Public License 1.1",
		"EUPL-1.2":             "European Union Public License 1.2",
		"CDDL-1.0":             "Common Development and Distribution License 1.0",
		"CDDL-1.1":             "Common Development and Distribution License 1.1",
		"EPL-1.0":              "Eclipse Public License 1.0",
		"EPL-2.0":              "Eclipse Public License 2.0",
		"MPL-1.1":              "Mozilla Public License 1.1",
		"MPL-2.0":              "Mozilla Public License 2.0",
		"OSL-3.0":              "Open Software License 3.0",
		"QPL-1.0":              "Q Public License 1.0",
		"Sleepycat":            "Sleepycat License",
	}

	return &LicenseAgent{
		highRiskLicenses: highRiskLicenses,
	}
}

// Name returns the identifier for this analysis agent.
func (la *LicenseAgent) Name() string {
	return "License Agent"
}

// Analyze examines the SBOM components for high-risk copyleft licenses.
// It returns a slice of AnalysisResult containing findings for components
// that use licenses identified as high-risk for compliance.
func (la *LicenseAgent) Analyze(ctx context.Context, sbom core.SBOM) ([]core.AnalysisResult, error) {
	var results []core.AnalysisResult

	for _, component := range sbom.Components {
		// Skip components without license information
		if component.License == "" {
			continue
		}

		// Check if the license is in our high-risk list
		if licenseDescription, isHighRisk := la.isHighRiskLicense(component.License); isHighRisk {
			// Determine severity based on license type
			severity := la.determineSeverity(component.License)
			
			// Create finding message
			finding := fmt.Sprintf("Component '%s' (v%s) uses high-risk copyleft license '%s' (%s). This may require source code disclosure or impose other compliance obligations.",
				component.Name,
				component.Version,
				component.License,
				licenseDescription)

			result := core.AnalysisResult{
				AgentName: la.Name(),
				Finding:   finding,
				Severity:  severity,
			}

			results = append(results, result)
		}
	}

	return results, nil
}

// isHighRiskLicense checks if a given license identifier is considered high-risk.
// It returns the license description and a boolean indicating if it's high-risk.
func (la *LicenseAgent) isHighRiskLicense(license string) (string, bool) {
	// Normalize the license string for comparison
	normalizedLicense := strings.TrimSpace(license)
	
	// Check exact match first
	if description, exists := la.highRiskLicenses[normalizedLicense]; exists {
		return description, true
	}
	
	// Check for common variations and shortened forms
	lowerLicense := strings.ToLower(normalizedLicense)
	for riskLicense, description := range la.highRiskLicenses {
		if strings.ToLower(riskLicense) == lowerLicense {
			return description, true
		}
		
		// Handle common shortened forms (e.g., "GPL-3.0" instead of "GPL-3.0-only")
		if strings.Contains(lowerLicense, "gpl") && strings.Contains(strings.ToLower(riskLicense), "gpl") {
			if extractVersionNumber(lowerLicense) == extractVersionNumber(strings.ToLower(riskLicense)) {
				return description, true
			}
		}
		
		if strings.Contains(lowerLicense, "agpl") && strings.Contains(strings.ToLower(riskLicense), "agpl") {
			if extractVersionNumber(lowerLicense) == extractVersionNumber(strings.ToLower(riskLicense)) {
				return description, true
			}
		}
		
		if strings.Contains(lowerLicense, "lgpl") && strings.Contains(strings.ToLower(riskLicense), "lgpl") {
			if extractVersionNumber(lowerLicense) == extractVersionNumber(strings.ToLower(riskLicense)) {
				return description, true
			}
		}
	}
	
	return "", false
}

// determineSeverity assigns a severity level based on the license type.
func (la *LicenseAgent) determineSeverity(license string) string {
	lowerLicense := strings.ToLower(license)
	
	// AGPL is considered the highest risk due to network copyleft provisions
	if strings.Contains(lowerLicense, "agpl") {
		return "Critical"
	}
	
	// Strong copyleft licenses (GPL)
	if strings.Contains(lowerLicense, "gpl") && !strings.Contains(lowerLicense, "lgpl") {
		return "High"
	}
	
	// Weaker copyleft licenses (LGPL, MPL, EPL, etc.)
	if strings.Contains(lowerLicense, "lgpl") || 
	   strings.Contains(lowerLicense, "mpl") || 
	   strings.Contains(lowerLicense, "epl") ||
	   strings.Contains(lowerLicense, "eupl") ||
	   strings.Contains(lowerLicense, "cddl") {
		return "Medium"
	}
	
	// Other copyleft licenses
	return "High"
}

// extractVersionNumber extracts version numbers from license strings for comparison.
func extractVersionNumber(license string) string {
	// Simple version extraction - looks for patterns like "2.0", "3.0", etc.
	if strings.Contains(license, "3.0") {
		return "3.0"
	}
	if strings.Contains(license, "2.1") {
		return "2.1"
	}
	if strings.Contains(license, "2.0") {
		return "2.0"
	}
	if strings.Contains(license, "1.1") {
		return "1.1"
	}
	if strings.Contains(license, "1.0") {
		return "1.0"
	}
	return ""
}