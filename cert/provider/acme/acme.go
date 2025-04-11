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

func (p *ACMEProvider) GetCertificate(host string, email string) (*tls.Certificate, error) {
	p.manager.HostPolicy = autocert.HostWhitelist(host)
	p.manager.Email = email

	cert, err := p.manager.GetCertificate(&tls.ClientHelloInfo{
		ServerName: host,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting certificate: %v", err)
	}

	return cert, nil
}

func (p *ACMEProvider) RenewCertificate(host string, email string) error {
	p.manager.Email = email
	return nil
}

func (p *ACMEProvider) ValidateCertificate(host string) bool {
	certFile := filepath.Join(p.CertDir, host)
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return false
	}
	return true
}
