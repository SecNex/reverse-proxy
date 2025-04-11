package provider

import (
	"crypto/tls"
)

// CertificateProvider definiert das Interface für Zertifikatsprovider
type CertificateProvider interface {
	// GetCertificate gibt das Zertifikat für den angegebenen Host zurück
	GetCertificate(host string) (*tls.Certificate, error)
	// RenewCertificate erneuert das Zertifikat für den angegebenen Host
	RenewCertificate(host string) error
	// ValidateCertificate prüft, ob das Zertifikat gültig ist
	ValidateCertificate(host string) bool
}

// BaseProvider enthält gemeinsame Felder für alle Provider
type BaseProvider struct {
	CertDir string
}
