// Package vectordb provides a simple in-memory vector database for storing and querying embeddings.
package vectordb

import (
	"fmt"
	"math"
	"sort"
)

// Document represents a document stored in the vector database.
type Document struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	Vector   []float64 `json:"vector"`
	Metadata map[string]interface{} `json:"metadata"`
}

// SearchResult represents a search result with similarity score.
type SearchResult struct {
	Document   Document `json:"document"`
	Similarity float64  `json:"similarity"`
}

// MemoryVectorDB is a simple in-memory vector database.
type MemoryVectorDB struct {
	documents map[string]Document
}

// NewMemoryVectorDB creates a new instance of MemoryVectorDB.
func NewMemoryVectorDB() *MemoryVectorDB {
	return &MemoryVectorDB{
		documents: make(map[string]Document),
	}
}

// Add adds a document to the vector database.
func (m *MemoryVectorDB) Add(doc Document) error {
	if doc.ID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}
	if len(doc.Vector) == 0 {
		return fmt.Errorf("document vector cannot be empty")
	}
	
	m.documents[doc.ID] = doc
	return nil
}

// Get retrieves a document by ID.
func (m *MemoryVectorDB) Get(id string) (Document, bool) {
	doc, exists := m.documents[id]
	return doc, exists
}

// Delete removes a document from the database.
func (m *MemoryVectorDB) Delete(id string) bool {
	if _, exists := m.documents[id]; exists {
		delete(m.documents, id)
		return true
	}
	return false
}

// Search performs similarity search and returns top k most similar documents.
func (m *MemoryVectorDB) Search(queryVector []float64, k int) ([]SearchResult, error) {
	if len(queryVector) == 0 {
		return nil, fmt.Errorf("query vector cannot be empty")
	}

	var results []SearchResult
	
	// Calculate cosine similarity for each document
	for _, doc := range m.documents {
		if len(doc.Vector) != len(queryVector) {
			continue // Skip documents with incompatible vector dimensions
		}
		
		similarity := cosineSimilarity(queryVector, doc.Vector)
		results = append(results, SearchResult{
			Document:   doc,
			Similarity: similarity,
		})
	}
	
	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
	
	// Return top k results
	if k > len(results) {
		k = len(results)
	}
	
	return results[:k], nil
}

// Size returns the number of documents in the database.
func (m *MemoryVectorDB) Size() int {
	return len(m.documents)
}

// Clear removes all documents from the database.
func (m *MemoryVectorDB) Clear() {
	m.documents = make(map[string]Document)
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}
	
	var dotProduct, normA, normB float64
	
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0.0 || normB == 0.0 {
		return 0.0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}