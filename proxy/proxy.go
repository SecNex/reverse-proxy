package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
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
	host := r.Host
	if strings.Contains(host, "%") {
		host = strings.Split(host, "%")[0]
	}

	if !isValidHost(host) {
		http.Error(w, "Ungültiger Host", http.StatusBadRequest)
		return
	}

	if host == "localhost" || host == "127.0.0.1" {
		http.ServeFile(w, r, "www/index.html")
		return
	}

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

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

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

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func isValidHost(host string) bool {
	if strings.Contains(host, ":") {
		ip := net.ParseIP(host)
		return ip != nil
	}

	if strings.Contains(host, ".") {
		return true
	}

	return false
}

func (rp *ReverseProxy) Start(port int, useSSL bool, certFile, keyFile string) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rp,
	}

	if useSSL {
		if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err == nil {
			if len(cert.Certificate) > 0 {
				if x509Cert, err := x509.ParseCertificate(cert.Certificate[0]); err == nil {
					if time.Now().Before(x509Cert.NotAfter) {
						server.TLSConfig = &tls.Config{
							Certificates: []tls.Certificate{cert},
						}
						return server.ListenAndServeTLS("", "")
					}
				}
			}
		}

		var err error
		certFile, keyFile, err = rp.certManager.GenerateSelfSignedCert("localhost", "ssl@example.local")
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
