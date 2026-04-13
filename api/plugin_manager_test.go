package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func resetPluginManager() {
	defaultPluginManager.mu.Lock()
	defer defaultPluginManager.mu.Unlock()
	defaultPluginManager.plugins = make(map[string]*Plugin)
}

func TestPluginManager_RegisterAndGetEnabled(t *testing.T) {
	resetPluginManager()

	defaultPluginManager.Register(&Plugin{Name: "alpha", URL: "http://alpha.example.com", Enabled: true})
	defaultPluginManager.Register(&Plugin{Name: "beta", URL: "http://beta.example.com", Enabled: false})

	enabled := defaultPluginManager.GetEnabled()
	if len(enabled) != 1 {
		t.Fatalf("expected 1 enabled plugin, got %d", len(enabled))
	}
	if enabled[0].Name != "alpha" {
		t.Errorf("expected alpha, got %s", enabled[0].Name)
	}
}

func TestPluginManager_Remove(t *testing.T) {
	resetPluginManager()

	defaultPluginManager.Register(&Plugin{Name: "gamma", URL: "http://gamma.example.com", Enabled: true})
	defaultPluginManager.Remove("gamma")

	all := defaultPluginManager.GetAll()
	if len(all) != 0 {
		t.Errorf("expected 0 plugins after removal, got %d", len(all))
	}
}

func TestPlugin_GetHTTPClient_DefaultTimeout(t *testing.T) {
	p := &Plugin{Name: "test", URL: "http://test.com", Timeout: 0}
	client := p.GetHTTPClient()
	if client.Timeout.Seconds() != 10 {
		t.Errorf("expected default timeout 10s, got %v", client.Timeout)
	}
}

func TestListPluginsHandler(t *testing.T) {
	resetPluginManager()
	defaultPluginManager.Register(&Plugin{Name: "p1", URL: "http://p1.com", Enabled: true})

	req := httptest.NewRequest(http.MethodGet, "/api/plugins", nil)
	w := httptest.NewRecorder()
	ListPluginsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var plugins []Plugin
	if err := json.NewDecoder(w.Body).Decode(&plugins); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(plugins) != 1 {
		t.Errorf("expected 1 plugin in response, got %d", len(plugins))
	}
}

func TestRegisterPluginHandler_Valid(t *testing.T) {
	resetPluginManager()

	body, _ := json.Marshal(Plugin{Name: "new", URL: "http://new.com", Enabled: true, Timeout: 5})
	req := httptest.NewRequest(http.MethodPost, "/api/plugins", bytes.NewReader(body))
	w := httptest.NewRecorder()
	RegisterPluginHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if len(defaultPluginManager.GetAll()) != 1 {
		t.Errorf("expected plugin to be registered")
	}
}

func TestRegisterPluginHandler_MissingFields(t *testing.T) {
	resetPluginManager()

	body, _ := json.Marshal(Plugin{Name: "", URL: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/plugins", bytes.NewReader(body))
	w := httptest.NewRecorder()
	RegisterPluginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
