package provider

import (
	"crypto/tls"
)

type CertificateProvider interface {
	GetCertificate(host string, email string) (*tls.Certificate, error)
	RenewCertificate(host string, email string) error
	ValidateCertificate(host string) bool
}

type BaseProvider struct {
	CertDir string
}
