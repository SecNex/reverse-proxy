package self

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/secnex/reverse-proxy/cert/provider"
)

type SelfSignedProvider struct {
	provider.BaseProvider
}

func NewProvider(certDir string) *SelfSignedProvider {
	return &SelfSignedProvider{
		BaseProvider: provider.BaseProvider{
			CertDir: certDir,
		},
	}
}

func (p *SelfSignedProvider) GetCertificate(host string) (*tls.Certificate, error) {
	certFile := filepath.Join(p.CertDir, fmt.Sprintf("%s.crt", host))
	keyFile := filepath.Join(p.CertDir, fmt.Sprintf("%s.key", host))

	// Prüfe, ob die Zertifikate bereits existieren
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		if err := p.generateCertificate(host); err != nil {
			return nil, err
		}
	}

	// Lade das Zertifikat
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (p *SelfSignedProvider) RenewCertificate(host string) error {
	return p.generateCertificate(host)
}

func (p *SelfSignedProvider) ValidateCertificate(host string) bool {
	certFile := filepath.Join(p.CertDir, fmt.Sprintf("%s.crt", host))
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (p *SelfSignedProvider) generateCertificate(host string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"SecNex Reverse Proxy"},
			OrganizationalUnit: []string{"Self Signed Certificate"},
			CommonName:         host,
			Country:            []string{"DE"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365), // 1 Jahr

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certFile := filepath.Join(p.CertDir, fmt.Sprintf("%s.crt", host))
	keyFile := filepath.Join(p.CertDir, fmt.Sprintf("%s.key", host))

	// Speichere das Zertifikat
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	certOut.Close()

	// Speichere den privaten Schlüssel
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	return nil
}
