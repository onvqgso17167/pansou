package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// sampleResults returns a slice of SearchResult for use in filter tests.
func sampleResults() []SearchResult {
	return []SearchResult{
		{Name: "Movie HD", URL: "https://pan.baidu.com/s/abc", Type: "baidu"},
		{Name: "Music Album", URL: "https://pan.quark.cn/s/xyz", Type: "quark"},
		{Name: "Document PDF", URL: "https://pan.baidu.com/s/def", Type: "baidu"},
		{Name: "Software Pack", URL: "https://aliyundrive.com/s/ghi", Type: "aliyun"},
		{Name: "", URL: "https://pan.baidu.com/s/empty", Type: "baidu"},
	}
}

// TestFilterResults_NoFilter verifies that without any filter all results are returned.
func TestFilterResults_NoFilter(t *testing.T) {
	results := sampleResults()
	filtered := filterResults(results, "", "")
	assert.Equal(t, len(results), len(filtered))
}

// TestFilterResults_ByType checks that filtering by type returns only matching entries.
func TestFilterResults_ByType(t *testing.T) {
	results := sampleResults()
	filtered := filterResults(results, "baidu", "")
	for _, r := range filtered {
		assert.Equal(t, "baidu", r.Type)
	}
	// Expect 3 baidu results (including the empty-name one)
	assert.Equal(t, 3, len(filtered))
}

// TestFilterResults_ByKeyword checks that keyword filtering returns only matching entries.
func TestFilterResults_ByKeyword(t *testing.T) {
	results := sampleResults()
	filtered := filterResults(results, "", "Movie")
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "Movie HD", filtered[0].Name)
}

// TestFilterResults_ByTypeAndKeyword checks combined type + keyword filtering.
func TestFilterResults_ByTypeAndKeyword(t *testing.T) {
	results := sampleResults()
	filtered := filterResults(results, "baidu", "Document")
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "Document PDF", filtered[0].Name)
}

// TestFilterResults_NoMatch ensures an empty slice is returned when nothing matches.
func TestFilterResults_NoMatch(t *testing.T) {
	results := sampleResults()
	filtered := filterResults(results, "nonexistent", "")
	assert.Empty(t, filtered)
}

// TestMatchFilter_CaseInsensitive verifies that keyword matching is case-insensitive.
func TestMatchFilter_CaseInsensitive(t *testing.T) {
	r := SearchResult{Name: "Hello World", URL: "https://example.com", Type: "test"}
	assert.True(t, matchFilter(r, "", "hello"))
	assert.True(t, matchFilter(r, "", "WORLD"))
	assert.False(t, matchFilter(r, "", "missing"))
}

// TestMatchFilter_TypeOnly checks type-only matching.
func TestMatchFilter_TypeOnly(t *testing.T) {
	r := SearchResult{Name: "Test", URL: "https://example.com", Type: "quark"}
	assert.True(t, matchFilter(r, "quark", ""))
	assert.False(t, matchFilter(r, "baidu", ""))
}

// TestFilterMergedByType verifies grouping of results by their type.
func TestFilterMergedByType(t *testing.T) {
	results := sampleResults()
	grouped := filterMergedByType(results)

	// baidu should have 3 entries
	assert.Equal(t, 3, len(grouped["baidu"]))
	// quark should have 1 entry
	assert.Equal(t, 1, len(grouped["quark"]))
	// aliyun should have 1 entry
	assert.Equal(t, 1, len(grouped["aliyun"]))
}
