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

	log.Println("Starting reverse proxy...")

	if err := os.MkdirAll(certDir+"/self", 0755); err != nil {
		log.Fatalf("Error creating self-signed certificate directory: %v", err)
	}
	if err := os.MkdirAll(certDir+"/acme", 0755); err != nil {
		log.Fatalf("Error creating ACME certificate directory: %v", err)
	}

	if err := os.MkdirAll(wwwDir, 0755); err != nil {
		log.Fatalf("Error creating www directory: %v", err)
	}

	dbManager, err := proxy.NewDBManager()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	certManager := cert.NewCertManager(certDir)
	configCache := proxy.NewConfigCache(dbManager, certManager)
	apiServer := server.NewAPIServer()
	reverseProxy := proxy.NewReverseProxy(configCache, certManager, apiServer)

	if err := configCache.LoadFromDB(); err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	go func() {
		if err := apiServer.Start(8081); err != nil {
			log.Printf("Error starting API server: %v", err)
		}
	}()

	localserverConfig := proxy.ProxyConfig{
		Protocol: "http",
		Host:     "localhost",
		Port:     8080,
		SSL:      true,
	}
	configCache.Set("localserver", localserverConfig)
	apiServer.SetActiveConfig("localserver", true)

	go func() {
		log.Println("Starting HTTP server on port 80...")
		if err := reverseProxy.Start(80, false, "", ""); err != nil {
			log.Printf("Fehler beim Starten des HTTP-Servers: %v", err)
		}
	}()

	log.Println("Starting HTTPS server on port 443...")
	if err := reverseProxy.Start(443, true, "", ""); err != nil {
		log.Fatalf("Fehler beim Starten des HTTPS-Servers: %v", err)
	}
}
