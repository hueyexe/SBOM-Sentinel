// Package ingestion provides CycloneDX JSON parsing functionality.
package ingestion

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hueyexe/SBOM-Sentinel/internal/core"
)

// CycloneDXParser implements the Parser interface for CycloneDX JSON format.
type CycloneDXParser struct{}

// NewCycloneDXParser creates a new instance of CycloneDXParser.
func NewCycloneDXParser() *CycloneDXParser {
	return &CycloneDXParser{}
}

// cycloneDXDocument represents the top-level structure of a CycloneDX JSON document.
type cycloneDXDocument struct {
	BOMFormat    string               `json:"bomFormat"`
	SpecVersion  string               `json:"specVersion"`
	SerialNumber string               `json:"serialNumber"`
	Version      int                  `json:"version"`
	Metadata     *cycloneDXMetadata   `json:"metadata,omitempty"`
	Components   []cycloneDXComponent `json:"components,omitempty"`
	Properties   []cycloneDXProperty  `json:"properties,omitempty"`
}

// cycloneDXMetadata represents the metadata section of a CycloneDX document.
type cycloneDXMetadata struct {
	Timestamp  string                  `json:"timestamp,omitempty"`
	Tools      []cycloneDXTool         `json:"tools,omitempty"`
	Authors    []cycloneDXOrganization `json:"authors,omitempty"`
	Component  *cycloneDXComponent     `json:"component,omitempty"`
	Supplier   *cycloneDXOrganization  `json:"supplier,omitempty"`
	Properties []cycloneDXProperty     `json:"properties,omitempty"`
}

// cycloneDXComponent represents a component in a CycloneDX document.
type cycloneDXComponent struct {
	Type       string                 `json:"type"`
	BOMRef     string                 `json:"bom-ref,omitempty"`
	Supplier   *cycloneDXOrganization `json:"supplier,omitempty"`
	Author     string                 `json:"author,omitempty"`
	Publisher  string                 `json:"publisher,omitempty"`
	Group      string                 `json:"group,omitempty"`
	Name       string                 `json:"name"`
	Version    string                 `json:"version"`
	PURL       string                 `json:"purl,omitempty"`
	Licenses   []cycloneDXLicense     `json:"licenses,omitempty"`
	Properties []cycloneDXProperty    `json:"properties,omitempty"`
}

// cycloneDXLicense represents a license in a CycloneDX document.
type cycloneDXLicense struct {
	License *cycloneDXLicenseChoice `json:"license,omitempty"`
}

// cycloneDXLicenseChoice represents the license choice structure.
type cycloneDXLicenseChoice struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}

// cycloneDXTool represents a tool in the metadata.
type cycloneDXTool struct {
	Vendor  string `json:"vendor,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// cycloneDXOrganization represents an organization in a CycloneDX document.
type cycloneDXOrganization struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// cycloneDXProperty represents a property in a CycloneDX document.
type cycloneDXProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Parse implements the Parser interface for CycloneDX JSON format.
// It reads a CycloneDX JSON document and converts it to our core SBOM model.
func (p *CycloneDXParser) Parse(r io.Reader) (*core.SBOM, error) {
	var doc cycloneDXDocument

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode CycloneDX JSON: %w", err)
	}

	// Validate that this is a CycloneDX document
	if doc.BOMFormat != "CycloneDX" {
		return nil, fmt.Errorf("invalid BOM format: expected 'CycloneDX', got '%s'", doc.BOMFormat)
	}

	// Convert to our core SBOM model
	sbom := &core.SBOM{
		ID:         doc.SerialNumber,
		Components: make([]core.Component, 0, len(doc.Components)),
		Metadata:   make(map[string]string),
	}

	// Set SBOM name from metadata if available
	if doc.Metadata != nil && doc.Metadata.Component != nil {
		sbom.Name = doc.Metadata.Component.Name
	}
	if sbom.Name == "" {
		sbom.Name = "Unnamed SBOM"
	}

	// Add metadata
	sbom.Metadata["bomFormat"] = doc.BOMFormat
	sbom.Metadata["specVersion"] = doc.SpecVersion
	if doc.Metadata != nil && doc.Metadata.Timestamp != "" {
		sbom.Metadata["timestamp"] = doc.Metadata.Timestamp
	}

	// Add properties as metadata
	for _, prop := range doc.Properties {
		sbom.Metadata[prop.Name] = prop.Value
	}

	// Convert components
	for _, comp := range doc.Components {
		component := core.Component{
			Name:    comp.Name,
			Version: comp.Version,
			PURL:    comp.PURL,
		}

		// Extract license information
		if len(comp.Licenses) > 0 && comp.Licenses[0].License != nil {
			license := comp.Licenses[0].License
			if license.ID != "" {
				component.License = license.ID
			} else if license.Name != "" {
				component.License = license.Name
			}
		}

		sbom.Components = append(sbom.Components, component)
	}

	return sbom, nil
}
