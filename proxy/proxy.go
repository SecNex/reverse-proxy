package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/secnex/reverse-proxy/cert"
	"github.com/secnex/reverse-proxy/server"
)

type ReverseProxy struct {
	configCache *ConfigCache
	certManager *cert.CertManager
	apiServer   *server.APIServer
}

func NewReverseProxy(configCache *ConfigCache, certManager *cert.CertManager, apiServer *server.APIServer) *ReverseProxy {
	return &ReverseProxy{
		configCache: configCache,
		certManager: certManager,
		apiServer:   apiServer,
	}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Hole den Host aus dem Request
	host := r.Host

	// Wenn der Host localhost oder 127.0.0.1 ist, serviere die statische index.html
	if host == "localhost" || host == "127.0.0.1" {
		http.ServeFile(w, r, "www/index.html")
		return
	}

	// Hole die Konfiguration für den Host
	config, exists := rp.configCache.Get(host)
	if !exists {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.ServeFile(w, r, "www/404.html")
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Host not found"}`))
		}
		return
	}

	// Prüfe, ob die Konfiguration aktiv ist
	if !rp.apiServer.IsActiveConfig(host) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.ServeFile(w, r, "www/503.html")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "Service Unavailable"}`))
		}
		return
	}

	// Erstelle einen neuen Request für den Ziel-Server
	req, err := http.NewRequest(r.Method, config.Protocol+"://"+config.Host+":"+strconv.Itoa(config.Port)+r.URL.Path, r.Body)
	if err != nil {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.ServeFile(w, r, "www/503.html")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "Service Unavailable"}`))
		}
		return
	}

	// Kopiere die Header vom Original-Request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Führe den Request aus
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.ServeFile(w, r, "www/502.html")
		} else {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Bad Gateway"}`))
		}
		return
	}
	defer resp.Body.Close()

	// Kopiere die Header vom Response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Kopiere den Status Code
	w.WriteHeader(resp.StatusCode)

	// Kopiere den Body
	io.Copy(w, resp.Body)
}

func (rp *ReverseProxy) Start(port int, useSSL bool, certFile, keyFile string) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rp,
	}

	if useSSL {
		// Prüfe, ob die Zertifikatsdateien existieren und gültig sind
		if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err == nil {
			// Prüfe die Gültigkeit des Zertifikats
			if len(cert.Certificate) > 0 {
				if x509Cert, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
					if time.Now().Before(x509Cert.NotAfter) {
						// Zertifikat ist gültig, verwende es
						server.TLSConfig = &tls.Config{
							Certificates: []tls.Certificate{cert},
						}
						return server.ListenAndServeTLS("", "")
					}
				}
			}
		}

		// Zertifikat existiert nicht oder ist ungültig, generiere ein neues
		var err error
		certFile, keyFile, err = rp.certManager.GenerateSelfSignedCert("localhost")
		if err != nil {
			return fmt.Errorf("error generating self-signed certificate: %v", err)
		}

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("error loading certificates: %v", err)
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		return server.ListenAndServeTLS("", "")
	}

	return server.ListenAndServe()
}
