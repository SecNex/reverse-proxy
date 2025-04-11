package acme

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"

	"github.com/secnex/reverse-proxy/cert/provider"
	"golang.org/x/crypto/acme/autocert"
)

type ACMEProvider struct {
	provider.BaseProvider
	manager *autocert.Manager
}

func NewProvider(certDir string) *ACMEProvider {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache(certDir),
		HostPolicy: autocert.HostWhitelist(),
		Prompt:     autocert.AcceptTOS,
	}

	return &ACMEProvider{
		BaseProvider: provider.BaseProvider{
			CertDir: certDir,
		},
		manager: manager,
	}
}

func (p *ACMEProvider) GetCertificate(host string) (*tls.Certificate, error) {
	// Füge den Host zur Whitelist hinzu
	p.manager.HostPolicy = autocert.HostWhitelist(host)

	// Hole das Zertifikat vom ACME-Server
	cert, err := p.manager.GetCertificate(&tls.ClientHelloInfo{
		ServerName: host,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting certificate: %v", err)
	}

	return cert, nil
}

func (p *ACMEProvider) RenewCertificate(host string) error {
	// ACME erneuert Zertifikate automatisch
	return nil
}

func (p *ACMEProvider) ValidateCertificate(host string) bool {
	certFile := filepath.Join(p.CertDir, host)
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return false
	}
	return true
}
