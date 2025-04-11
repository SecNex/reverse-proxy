package main

import (
	"log"
	"os"

	"github.com/secnex/reverse-proxy/cert"
	"github.com/secnex/reverse-proxy/proxy"
	"github.com/secnex/reverse-proxy/server"
)

func main() {
	certDir := "certs"
	wwwDir := "www"
	// Erstelle die Verzeichnisse für die Zertifikate
	if err := os.MkdirAll(certDir+"/self", 0755); err != nil {
		log.Fatalf("Error creating certificate directory: %v", err)
	}
	if err := os.MkdirAll(certDir+"/acme", 0755); err != nil {
		log.Fatalf("Error creating ACME directory: %v", err)
	}

	// Erstelle die Verzeichnisse für die statischen Dateien
	if err := os.MkdirAll(wwwDir, 0755); err != nil {
		log.Fatalf("Error creating www directory: %v", err)
	}

	// Initialisiere die Komponenten
	configCache := proxy.NewConfigCache()
	certManager := cert.NewCertManager(certDir)
	apiServer := server.NewAPIServer()
	reverseProxy := proxy.NewReverseProxy(configCache, certManager, apiServer)

	// Starte den API-Server
	go func() {
		if err := apiServer.Start(8081); err != nil {
			log.Printf("Error starting API server: %v", err)
		}
	}()

	// Konfiguriere localhost
	// localhostConfig := proxy.ProxyConfig{
	// 	TargetURL: "http://localhost:8082",
	// 	SSL:       true,
	// }
	// configCache.Set("localhost", localhostConfig)
	// apiServer.SetActiveConfig("localhost", true)

	// Konfiguriere localserver
	localserverConfig := proxy.ProxyConfig{
		Protocol: "http",
		Host:     "localhost",
		Port:     8080,
		SSL:      false,
	}
	configCache.Set("localserver", localserverConfig)
	apiServer.SetActiveConfig("localserver", true)

	// Starte den Reverse Proxy für HTTP (Port 80)
	go func() {
		log.Println("Starting HTTP server on port 80...")
		if err := reverseProxy.Start(80, false, "", ""); err != nil {
			log.Printf("Error starting HTTP server: %v", err)
		}
	}()

	// Starte den Reverse Proxy für HTTPS (Port 443)
	log.Println("Starting HTTPS server on port 443...")
	if err := reverseProxy.Start(443, true, "", ""); err != nil {
		log.Fatalf("Error starting HTTPS server: %v", err)
	}
}
