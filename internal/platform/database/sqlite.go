// Package database provides concrete implementations of storage interfaces using SQLite.
package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/chrisclapham/SBOM-Sentinel/internal/core"
	"github.com/chrisclapham/SBOM-Sentinel/internal/platform/storage"
)

// SQLiteRepository implements the storage.Repository interface using SQLite.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository instance.
// It initializes the database connection and creates the necessary tables.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the necessary tables for storing SBOM data.
func (r *SQLiteRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sboms (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		components TEXT NOT NULL, -- JSON-encoded components
		metadata TEXT NOT NULL,   -- JSON-encoded metadata
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_sboms_name ON sboms(name);
	CREATE INDEX IF NOT EXISTS idx_sboms_created_at ON sboms(created_at);
	`

	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Store persists an SBOM document to the SQLite database.
func (r *SQLiteRepository) Store(ctx context.Context, sbom core.SBOM) error {
	// Serialize components to JSON
	componentsJSON, err := json.Marshal(sbom.Components)
	if err != nil {
		return fmt.Errorf("failed to marshal components: %w", err)
	}

	// Serialize metadata to JSON
	metadataJSON, err := json.Marshal(sbom.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	now := time.Now()

	// Check if SBOM already exists
	var existingID string
	err = r.db.QueryRowContext(ctx, "SELECT id FROM sboms WHERE id = ?", sbom.ID).Scan(&existingID)
	
	if err == sql.ErrNoRows {
		// Insert new SBOM
		query := `
			INSERT INTO sboms (id, name, components, metadata, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = r.db.ExecContext(ctx, query, sbom.ID, sbom.Name, string(componentsJSON), string(metadataJSON), now, now)
		if err != nil {
			return fmt.Errorf("failed to insert SBOM: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing SBOM: %w", err)
	} else {
		// Update existing SBOM
		query := `
			UPDATE sboms 
			SET name = ?, components = ?, metadata = ?, updated_at = ?
			WHERE id = ?
		`
		_, err = r.db.ExecContext(ctx, query, sbom.Name, string(componentsJSON), string(metadataJSON), now, sbom.ID)
		if err != nil {
			return fmt.Errorf("failed to update SBOM: %w", err)
		}
	}

	return nil
}

// FindByID retrieves an SBOM document by its unique identifier.
func (r *SQLiteRepository) FindByID(ctx context.Context, id string) (*core.SBOM, error) {
	query := `
		SELECT id, name, components, metadata, created_at, updated_at
		FROM sboms
		WHERE id = ?
	`

	var sbom core.SBOM
	var componentsJSON, metadataJSON string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sbom.ID,
		&sbom.Name,
		&componentsJSON,
		&metadataJSON,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // SBOM not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query SBOM: %w", err)
	}

	// Deserialize components from JSON
	if err := json.Unmarshal([]byte(componentsJSON), &sbom.Components); err != nil {
		return nil, fmt.Errorf("failed to unmarshal components: %w", err)
	}

	// Deserialize metadata from JSON
	if err := json.Unmarshal([]byte(metadataJSON), &sbom.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &sbom, nil
}

// Close closes the database connection.
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Verify that SQLiteRepository implements the storage.Repository interface.
var _ storage.Repository = (*SQLiteRepository)(nil)