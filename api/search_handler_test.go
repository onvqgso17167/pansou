package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupSearchRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/search", SearchHandler)
	return r
}

func TestSearchHandler_MissingKeyword(t *testing.T) {
	r := setupSearchRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSearchHandler_WithKeyword(t *testing.T) {
	r := setupSearchRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/search?k=golang", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Keyword != "golang" {
		t.Errorf("expected keyword 'golang', got '%s'", resp.Keyword)
	}
}

func TestParseResultsFromJSON_Valid(t *testing.T) {
	data := []byte(`[{"name":"test","url":"https://pan.baidu.com/s/abc","type":"baidu","source":"chan1"}]`)
	results, err := ParseResultsFromJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "test" {
		t.Errorf("expected name 'test', got '%s'", results[0].Name)
	}
}

func TestParseResultsFromJSON_Invalid(t *testing.T) {
	_, err := ParseResultsFromJSON([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
