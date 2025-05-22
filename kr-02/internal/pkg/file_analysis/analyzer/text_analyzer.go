package analyzer

import (
	"regexp"
	"strings"
	"unicode"
)

// TextAnalyzer provides methods for analyzing text content
type TextAnalyzer struct {
	// Common words to ignore in analysis (stop words)
	StopWords map[string]bool
}

// NewTextAnalyzer creates a new TextAnalyzer instance
func NewTextAnalyzer() *TextAnalyzer {
	// Initialize with common English stop words
	stopWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "with": true,
		"by": true, "of": true, "about": true, "from": true,
		"this": true, "that": true, "these": true, "those": true,
		"it": true, "its": true, "it's": true, "they": true, "them": true, "their": true,
	}

	return &TextAnalyzer{
		StopWords: stopWords,
	}
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

// GetSignificantWords returns words after removing stop words and punctuation
func (a *TextAnalyzer) GetSignificantWords(content string) []string {
	// Convert to lowercase
	content = strings.ToLower(content)

	// Remove punctuation
	var sb strings.Builder
	for _, r := range content {
		if !unicode.IsPunct(r) {
			sb.WriteRune(r)
		} else {
			sb.WriteRune(' ')
		}
	}
	content = sb.String()

	// Split into words
	words := strings.Fields(content)

	// Filter out stop words
	var significantWords []string
	for _, word := range words {
		if !a.StopWords[word] {
			significantWords = append(significantWords, word)
		}
	}

	return significantWords
}

// GetNGrams returns n-grams (sequences of n consecutive words) from the content
func (a *TextAnalyzer) GetNGrams(content string, n int) []string {
	words := a.GetWords(content)
	if len(words) < n {
		return []string{}
	}

	var ngrams []string
	for i := 0; i <= len(words)-n; i++ {
		ngram := strings.Join(words[i:i+n], " ")
		ngrams = append(ngrams, ngram)
	}

	return ngrams
}

// RemoveExcessWhitespace normalizes whitespace in text
func (a *TextAnalyzer) RemoveExcessWhitespace(text string) string {
	// Replace multiple spaces, tabs, and newlines with a single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}
