package analysis

import (
	"context"
	"testing"

	"github.com/hueyexe/SBOM-Sentinel/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestLicenseAgent_Name(t *testing.T) {
	agent := NewLicenseAgent()
	assert.Equal(t, "License Agent", agent.Name())
}

func TestLicenseAgent_Analyze(t *testing.T) {
	tests := []struct {
		name               string
		sbom               core.SBOM
		expectedCount      int
		expectedFindings   []string
		expectedSeverities []string
	}{
		{
			name: "AGPL license detected - Critical severity",
			sbom: core.SBOM{
				ID:   "test-1",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "test-component",
						Version: "1.0.0",
						License: "AGPL-3.0-only",
					},
				},
			},
			expectedCount:      1,
			expectedFindings:   []string{"Component 'test-component' (v1.0.0) uses high-risk copyleft license 'AGPL-3.0-only'"},
			expectedSeverities: []string{"Critical"},
		},
		{
			name: "GPL license detected - High severity",
			sbom: core.SBOM{
				ID:   "test-2",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "gpl-component",
						Version: "2.1.0",
						License: "GPL-3.0-only",
					},
				},
			},
			expectedCount:      1,
			expectedFindings:   []string{"Component 'gpl-component' (v2.1.0) uses high-risk copyleft license 'GPL-3.0-only'"},
			expectedSeverities: []string{"High"},
		},
		{
			name: "LGPL license detected - Medium severity",
			sbom: core.SBOM{
				ID:   "test-3",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "lgpl-component",
						Version: "1.5.0",
						License: "LGPL-3.0-only",
					},
				},
			},
			expectedCount:      1,
			expectedFindings:   []string{"Component 'lgpl-component' (v1.5.0) uses high-risk copyleft license 'LGPL-3.0-only'"},
			expectedSeverities: []string{"Medium"},
		},
		{
			name: "Multiple high-risk licenses",
			sbom: core.SBOM{
				ID:   "test-4",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "component1",
						Version: "1.0.0",
						License: "AGPL-3.0-only",
					},
					{
						Name:    "component2",
						Version: "2.0.0",
						License: "GPL-2.0-only",
					},
					{
						Name:    "component3",
						Version: "1.1.0",
						License: "MPL-2.0",
					},
				},
			},
			expectedCount:      3,
			expectedFindings:   []string{"AGPL-3.0-only", "GPL-2.0-only", "MPL-2.0"},
			expectedSeverities: []string{"Critical", "High", "Medium"},
		},
		{
			name: "Safe licenses - no findings",
			sbom: core.SBOM{
				ID:   "test-5",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "safe-component1",
						Version: "1.0.0",
						License: "MIT",
					},
					{
						Name:    "safe-component2",
						Version: "2.0.0",
						License: "Apache-2.0",
					},
					{
						Name:    "safe-component3",
						Version: "1.5.0",
						License: "BSD-3-Clause",
					},
				},
			},
			expectedCount:      0,
			expectedFindings:   []string{},
			expectedSeverities: []string{},
		},
		{
			name: "Components without license information",
			sbom: core.SBOM{
				ID:   "test-6",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "no-license-component",
						Version: "1.0.0",
						License: "",
					},
					{
						Name:    "risky-component",
						Version: "2.0.0",
						License: "GPL-3.0-only",
					},
				},
			},
			expectedCount:      1,
			expectedFindings:   []string{"GPL-3.0-only"},
			expectedSeverities: []string{"High"},
		},
		{
			name: "License case variations",
			sbom: core.SBOM{
				ID:   "test-7",
				Name: "Test SBOM",
				Components: []core.Component{
					{
						Name:    "case-component1",
						Version: "1.0.0",
						License: "gpl-3.0",
					},
					{
						Name:    "case-component2",
						Version: "2.0.0",
						License: "LGPL-2.1",
					},
				},
			},
			expectedCount:      2,
			expectedFindings:   []string{"gpl-3.0", "LGPL-2.1"},
			expectedSeverities: []string{"High", "Medium"},
		},
		{
			name: "Empty SBOM",
			sbom: core.SBOM{
				ID:         "test-8",
				Name:       "Empty SBOM",
				Components: []core.Component{},
			},
			expectedCount:      0,
			expectedFindings:   []string{},
			expectedSeverities: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewLicenseAgent()
			ctx := context.Background()

			results, err := agent.Analyze(ctx, tt.sbom)

			// Check that no error occurred
			assert.NoError(t, err)

			// Check the number of results
			assert.Equal(t, tt.expectedCount, len(results))

			// Check each result
			for i, result := range results {
				assert.Equal(t, "License Agent", result.AgentName)

				if i < len(tt.expectedSeverities) {
					assert.Equal(t, tt.expectedSeverities[i], result.Severity)
				}

				if i < len(tt.expectedFindings) {
					assert.Contains(t, result.Finding, tt.expectedFindings[i])
				}
			}
		})
	}
}

func TestLicenseAgent_isHighRiskLicense(t *testing.T) {
	agent := NewLicenseAgent()

	tests := []struct {
		name        string
		license     string
		expectRisk  bool
		description string
	}{
		{
			name:        "AGPL-3.0-only exact match",
			license:     "AGPL-3.0-only",
			expectRisk:  true,
			description: "GNU Affero General Public License v3.0 only",
		},
		{
			name:        "GPL-3.0 shortened form",
			license:     "GPL-3.0",
			expectRisk:  true,
			description: "",
		},
		{
			name:       "MIT license - safe",
			license:    "MIT",
			expectRisk: false,
		},
		{
			name:       "Apache-2.0 license - safe",
			license:    "Apache-2.0",
			expectRisk: false,
		},
		{
			name:       "Empty license",
			license:    "",
			expectRisk: false,
		},
		{
			name:        "Case insensitive match",
			license:     "gpl-2.0-only",
			expectRisk:  true,
			description: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description, isRisk := agent.isHighRiskLicense(tt.license)
			assert.Equal(t, tt.expectRisk, isRisk)

			if tt.description != "" {
				assert.Equal(t, tt.description, description)
			}
		})
	}
}

func TestLicenseAgent_determineSeverity(t *testing.T) {
	agent := NewLicenseAgent()

	tests := []struct {
		name             string
		license          string
		expectedSeverity string
	}{
		{
			name:             "AGPL license - Critical",
			license:          "AGPL-3.0-only",
			expectedSeverity: "Critical",
		},
		{
			name:             "GPL license - High",
			license:          "GPL-3.0-only",
			expectedSeverity: "High",
		},
		{
			name:             "LGPL license - Medium",
			license:          "LGPL-2.1-only",
			expectedSeverity: "Medium",
		},
		{
			name:             "MPL license - Medium",
			license:          "MPL-2.0",
			expectedSeverity: "Medium",
		},
		{
			name:             "EPL license - Medium",
			license:          "EPL-2.0",
			expectedSeverity: "Medium",
		},
		{
			name:             "Unknown copyleft license - High",
			license:          "OSL-3.0",
			expectedSeverity: "High",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			severity := agent.determineSeverity(tt.license)
			assert.Equal(t, tt.expectedSeverity, severity)
		})
	}
}

func TestLicenseAgent_extractVersionNumber(t *testing.T) {
	tests := []struct {
		name            string
		license         string
		expectedVersion string
	}{
		{
			name:            "GPL 3.0",
			license:         "gpl-3.0-only",
			expectedVersion: "3.0",
		},
		{
			name:            "LGPL 2.1",
			license:         "lgpl-2.1-or-later",
			expectedVersion: "2.1",
		},
		{
			name:            "GPL 2.0",
			license:         "gpl-2.0",
			expectedVersion: "2.0",
		},
		{
			name:            "No version",
			license:         "MIT",
			expectedVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := extractVersionNumber(tt.license)
			assert.Equal(t, tt.expectedVersion, version)
		})
	}
}
