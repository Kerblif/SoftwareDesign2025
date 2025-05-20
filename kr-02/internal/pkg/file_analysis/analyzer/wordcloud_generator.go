package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/google/uuid"
)

// WordCloudGenerator provides methods for generating word clouds
type WordCloudGenerator struct {
	apiURL string
}

// NewWordCloudGenerator creates a new WordCloudGenerator instance
func NewWordCloudGenerator(apiURL string) *WordCloudGenerator {
	if apiURL == "" {
		apiURL = "https://quickchart.io/wordcloud"
	}
	return &WordCloudGenerator{
		apiURL: apiURL,
	}
}

// WordItem represents a word and its frequency for the word cloud API
type WordItem struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

// GenerateWordCloud generates a word cloud image from the given words
func (g *WordCloudGenerator) GenerateWordCloud(ctx context.Context, words []string) ([]byte, string, error) {
	// Count word frequencies
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}
	
	// Convert to the format expected by the API
	var wordItems []WordItem
	for word, freq := range wordFreq {
		wordItems = append(wordItems, WordItem{
			Text:  word,
			Value: freq,
		})
	}
	
	// Prepare the request payload
	requestData := struct {
		Width  int        `json:"width"`
		Height int        `json:"height"`
		Words  []WordItem `json:"words"`
	}{
		Width:  800,
		Height: 400,
		Words:  wordItems,
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request data: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", g.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}
	
	// Read the response body (image data)
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Generate a unique location for the image
	location := uuid.New().String() + ".png"
	
	return imageData, location, nil
}