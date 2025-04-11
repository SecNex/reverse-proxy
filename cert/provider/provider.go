package provider

import (
	"crypto/tls"
)

type CertificateProvider interface {
	GetCertificate(host string) (*tls.Certificate, error)
	RenewCertificate(host string) error
	ValidateCertificate(host string) bool
}

type BaseProvider struct {
	CertDir string
}
