package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SearchResult represents a single search result item
type SearchResult struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Size    string `json:"size,omitempty"`
	Source  string `json:"source"`
	AddedAt string `json:"added_at,omitempty"`
}

// SearchResponse wraps the search results and metadata
type SearchResponse struct {
	Keyword string         `json:"keyword"`
	Total   int            `json:"total"`
	Results []SearchResult `json:"results"`
	Elapsed string         `json:"elapsed"`
}

// SearchHandler handles search requests
func SearchHandler(c *gin.Context) {
	start := time.Now()

	keyword := strings.TrimSpace(c.Query("k"))
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "keyword is required",
		})
		return
	}

	typeFilter := c.Query("t")

	// Fetch raw results from channels
	rawResults := fetchFromChannels(keyword)

	// Apply filters
	filtered := applyResultFilter(rawResults, typeFilter)

	resp := SearchResponse{
		Keyword: keyword,
		Total:   len(filtered),
		Results: filtered,
		Elapsed: time.Since(start).String(),
	}

	c.JSON(http.StatusOK, resp)
}

// fetchFromChannels simulates fetching results from multiple pan search channels
func fetchFromChannels(keyword string) []SearchResult {
	// Placeholder: in production this would fan-out to multiple sources
	_ = keyword
	return []SearchResult{}
}

// ParseResultsFromJSON parses a JSON byte slice into a slice of SearchResult
func ParseResultsFromJSON(data []byte) ([]SearchResult, error) {
	var results []SearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}
