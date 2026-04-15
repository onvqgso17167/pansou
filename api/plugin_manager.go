package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// Plugin represents a search channel plugin configuration
type Plugin struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
	Timeout int    `json:"timeout"` // seconds
}

// PluginManager manages the list of registered search plugins
type PluginManager struct {
	mu      sync.RWMutex
	plugins map[string]*Plugin
}

var defaultPluginManager = &PluginManager{
	plugins: make(map[string]*Plugin),
}

// Register adds or updates a plugin in the manager
func (pm *PluginManager) Register(p *Plugin) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[p.Name] = p
	log.Printf("[plugin] registered: %s (enabled=%v)", p.Name, p.Enabled)
}

// GetEnabled returns all enabled plugins
func (pm *PluginManager) GetEnabled() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	var result []*Plugin
	for _, p := range pm.plugins {
		if p.Enabled {
			result = append(result, p)
		}
	}
	return result
}

// GetAll returns all registered plugins
func (pm *PluginManager) GetAll() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	var result []*Plugin
	for _, p := range pm.plugins {
		result = append(result, p)
	}
	return result
}

// Remove deletes a plugin by name
func (pm *PluginManager) Remove(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.plugins, name)
}

// GetHTTPClient returns an http.Client with the plugin's configured timeout.
// Default timeout is increased to 15s to reduce premature timeouts on slow sources.
func (p *Plugin) GetHTTPClient() *http.Client {
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 15
	}
	return &http.Client{Timeout: time.Duration(timeout) * time.Second}
}

// ListPluginsHandler handles GET /api/plugins
func ListPluginsHandler(w http.ResponseWriter, r *http.Request) {
	plugins := defaultPluginManager.GetAll()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(plugins); err != nil {
		http.Error(w, "failed to encode plugins", http.StatusInternalServerError)
	}
}

// RegisterPluginHandler handles POST /api/plugins
func RegisterPluginHandler(w http.ResponseWriter, r *http.Request) {
	var p Plugin
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if p.Name == "" || p.URL == "" {
		http.Error(w, "name and url are required", http.StatusBadRequest)
		return
	}
	defaultPluginManager.Register(&p)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
