package cert

import (
	"crypto/tls"
	"fmt"
	"log"
	"path/filepath"

	"github.com/secnex/reverse-proxy/cert/provider"
	"github.com/secnex/reverse-proxy/cert/provider/acme"
	"github.com/secnex/reverse-proxy/cert/provider/self"
)

type CertManager struct {
	certDir   string
	providers map[string]provider.CertificateProvider
}

func NewCertManager(certDir string) *CertManager {
	cm := &CertManager{
		certDir:   certDir,
		providers: make(map[string]provider.CertificateProvider),
	}

	cm.providers["self"] = self.NewProvider(filepath.Join(certDir, "self"))
	cm.providers["acme"] = acme.NewProvider(filepath.Join(certDir, "acme"))

	return cm
}

func (cm *CertManager) GetCertificate(host string, providerType string) (*tls.Certificate, error) {
	provider, exists := cm.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}

	return provider.GetCertificate(host)
}

func (cm *CertManager) RenewCertificate(host string, providerType string) error {
	provider, exists := cm.providers[providerType]
	if !exists {
		return fmt.Errorf("provider %s not found", providerType)
	}

	return provider.RenewCertificate(host)
}

func (cm *CertManager) ValidateCertificate(host string, providerType string) bool {
	provider, exists := cm.providers[providerType]
	if !exists {
		return false
	}

	return provider.ValidateCertificate(host)
}

func (cm *CertManager) GenerateSelfSignedCert(host string) (string, string, error) {
	log.Printf("Generating self-signed certificate for %s...", host)
	_, err := cm.GetCertificate(host, "self")
	if err != nil {
		return "", "", err
	}

	certFile := filepath.Join(cm.certDir, "self", fmt.Sprintf("%s.crt", host))
	keyFile := filepath.Join(cm.certDir, "self", fmt.Sprintf("%s.key", host))
	log.Printf("Certificate generated for %s!", host)
	return certFile, keyFile, nil
}

func (cm *CertManager) LoadCert(certFile, keyFile string) (*tls.Certificate, error) {
	log.Printf("Loading certificate from %s and key from %s...", certFile, keyFile)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	log.Printf("Certificate loaded successfully for %s!", certFile)
	return &cert, nil
}
