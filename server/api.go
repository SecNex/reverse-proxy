package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type APIServer struct {
	activeConfigs map[string]bool
	mu            sync.RWMutex
}

func NewAPIServer() *APIServer {
	return &APIServer{
		activeConfigs: make(map[string]bool),
	}
}

func (s *APIServer) Start(port int) error {
	http.HandleFunc("/api/status", s.handleStatus)
	http.HandleFunc("/api/refresh", s.handleRefresh)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	activeSites := make([]string, 0, len(s.activeConfigs))
	for site := range s.activeConfigs {
		activeSites = append(activeSites, site)
	}

	response := struct {
		ActiveSites []string `json:"active_sites"`
	}{
		ActiveSites: activeSites,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Refresh triggered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) SetActiveConfig(site string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeConfigs[site] = active
}

func (s *APIServer) IsActiveConfig(site string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeConfigs[site]
}
