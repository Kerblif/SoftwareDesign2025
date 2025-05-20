package analyzer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// PlagiarismChecker provides methods for checking plagiarism
type PlagiarismChecker struct{}

// NewPlagiarismChecker creates a new PlagiarismChecker instance
func NewPlagiarismChecker() *PlagiarismChecker {
	return &PlagiarismChecker{}
}

// CheckPlagiarism checks if the content is plagiarized from any of the provided contents
// Returns true if plagiarism is detected, along with the IDs of similar files
func (c *PlagiarismChecker) CheckPlagiarism(ctx context.Context, content string, otherContents map[string]string) (bool, []string) {
	// For simplicity, we're just checking for exact matches
	// In a real-world scenario, you would use more sophisticated algorithms
	
	// Calculate hash of the current content
	currentHash := c.calculateHash(content)
	
	var similarFileIDs []string
	
	// Compare with other contents
	for fileID, otherContent := range otherContents {
		otherHash := c.calculateHash(otherContent)
		
		// If hashes match, it's a potential plagiarism
		if currentHash == otherHash {
			similarFileIDs = append(similarFileIDs, fileID)
		}
	}
	
	return len(similarFileIDs) > 0, similarFileIDs
}

// calculateHash calculates a SHA-256 hash of the content
func (c *PlagiarismChecker) calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}