// Package analysis provides interfaces and implementations for analyzing SBOM documents.
// This package supports a pluggable architecture where different analysis agents
// can perform specialized security, compliance, and health checks.
package analysis

import (
	"context"

	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
)

// AnalysisAgent defines the contract for analyzing SBOM documents.
// Each implementation focuses on a specific type of analysis such as
// vulnerability scanning, license compliance, or dependency health checks.
type AnalysisAgent interface {
	// Name returns a unique identifier for this analysis agent.
	// This name is used in AnalysisResult to identify the source of findings.
	Name() string
	
	// Analyze performs analysis on the provided SBOM and returns any findings.
	// The context can be used for cancellation and timeout control.
	// Returns a slice of AnalysisResult containing all findings, or an error
	// if the analysis cannot be completed.
	Analyze(ctx context.Context, sbom core.SBOM) ([]core.AnalysisResult, error)
}