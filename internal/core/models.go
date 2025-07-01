// Package core contains the domain models and business logic for SBOM Sentinel.
// This package has no external dependencies and represents the core of our hexagonal architecture.
package core

// Component represents a software component within an SBOM.
// It contains essential metadata about a software package or library.
type Component struct {
	// Name is the human-readable name of the component
	Name string `json:"name"`
	
	// Version is the version identifier of the component
	Version string `json:"version"`
	
	// PURL (Package URL) is a standardized way to identify and locate software packages
	PURL string `json:"purl"`
	
	// License is the license identifier or expression for the component
	License string `json:"license"`
}

// SBOM represents a Software Bill of Materials document.
// It contains a collection of components and associated metadata.
type SBOM struct {
	// ID is a unique identifier for this SBOM
	ID string `json:"id"`
	
	// Name is a human-readable name for this SBOM
	Name string `json:"name"`
	
	// Components is a slice of all software components included in this SBOM
	Components []Component `json:"components"`
	
	// Metadata contains additional key-value pairs of information about the SBOM
	Metadata map[string]string `json:"metadata"`
}

// AnalysisResult represents the outcome of running an analysis agent on an SBOM.
// It contains the findings and severity assessment from a specific analysis.
type AnalysisResult struct {
	// AgentName identifies which analysis agent produced this result
	AgentName string `json:"agent_name"`
	
	// Finding describes what was discovered during the analysis
	Finding string `json:"finding"`
	
	// Severity indicates the severity level of the finding (e.g., "low", "medium", "high", "critical")
	Severity string `json:"severity"`
}