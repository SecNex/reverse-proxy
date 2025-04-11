package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type APIServer struct {
	activeConfigs map[string]bool
	mu            sync.RWMutex
	rateLimiter   map[string]time.Time
	rateLimit     time.Duration
}

func NewAPIServer() *APIServer {
	return &APIServer{
		activeConfigs: make(map[string]bool),
		rateLimiter:   make(map[string]time.Time),
		rateLimit:     time.Second * 1,
	}
}

func (s *APIServer) Start(port int) error {
	http.HandleFunc("/api/status", s.handleStatus)
	http.HandleFunc("/api/refresh", s.handleRefresh)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if !s.checkRateLimit(r) {
		http.Error(w, "Zu viele Anfragen", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if !s.checkRateLimit(r) {
		http.Error(w, "Zu viele Anfragen", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Aktualisierung erfolgreich ausgelöst",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) checkRateLimit(r *http.Request) bool {
	ip := r.RemoteAddr
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	if lastRequest, exists := s.rateLimiter[ip]; exists {
		if now.Sub(lastRequest) < s.rateLimit {
			return false
		}
	}

	s.rateLimiter[ip] = now
	return true
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
