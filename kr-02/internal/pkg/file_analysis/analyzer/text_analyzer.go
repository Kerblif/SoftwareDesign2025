package analyzer

import (
	"strings"
)

// TextAnalyzer provides methods for analyzing text content
type TextAnalyzer struct{}

// NewTextAnalyzer creates a new TextAnalyzer instance
func NewTextAnalyzer() *TextAnalyzer {
	return &TextAnalyzer{}
}

// AnalyzeText analyzes text content and returns statistics
func (a *TextAnalyzer) AnalyzeText(content string) (paragraphCount, wordCount, characterCount int32) {
	// Count paragraphs (separated by double newlines)
	paragraphs := strings.Split(content, "\n\n")
	// Filter out empty paragraphs
	var nonEmptyParagraphs []string
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			nonEmptyParagraphs = append(nonEmptyParagraphs, p)
		}
	}
	paragraphCount = int32(len(nonEmptyParagraphs))
	
	// Count words
	words := strings.Fields(content)
	wordCount = int32(len(words))
	
	// Count characters (including whitespace)
	characterCount = int32(len(content))
	
	return paragraphCount, wordCount, characterCount
}

// GetWords returns a slice of all words in the content
func (a *TextAnalyzer) GetWords(content string) []string {
	return strings.Fields(content)
}