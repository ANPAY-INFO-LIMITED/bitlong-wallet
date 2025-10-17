package crt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/pkg/errors"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	country            = "CN"
	organization       = "btl"
	organizationalUnit = "IT"
	locality           = "Chengdu"
	province           = "Sichuan"
	commonName         = "localhost"

	day  = 24 * time.Hour
	year = 365 * day

	CertPath = "bitlong_pc/etc/tls/tls.cert"
	KeyPath  = "bitlong_pc/etc/tls/tls.key"
)

var (
	certNotExist = errors.New("TLS certificate does not exist")
	KeyNotExist  = errors.New("TLS private key does not exist")
)

func GenerateSelfSignedTlsCert() (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "ecdsa.GenerateKey")
	}

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
			Locality:           []string{locality},
			Province:           []string{province},
			CommonName:         commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(10 * year),

		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "x509.CreateCertificate")
	}

	privateKeyBits, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "x509.MarshalECPrivateKey")
	}

	certPEM := pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: certDER},
	)
	keyPEM := pem.EncodeToMemory(
		&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBits},
	)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "os.UserHomeDir")
	}
	certAbsolutePath := filepath.Join(homeDir, CertPath)
	keyAbsolutePath := filepath.Join(homeDir, KeyPath)

	certDir := filepath.Dir(certAbsolutePath)
	err = os.MkdirAll(certDir, 0644)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "os.MkdirAll certDir")
	}
	keyDir := filepath.Dir(keyAbsolutePath)
	err = os.MkdirAll(keyDir, 0644)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "os.MkdirAll keyDir")
	}

	if err = os.WriteFile(certAbsolutePath, certPEM, 0644); err != nil {
		return tls.Certificate{}, errors.Wrap(err, "os.WriteFile certPEM")
	}
	if err = os.WriteFile(keyAbsolutePath, keyPEM, 0600); err != nil {
		return tls.Certificate{}, errors.Wrap(err, "os.WriteFile keyPEM")
	}

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "tls.X509KeyPair")
	}

	return tlsCert, nil
}

func CheckCertExist() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "os.UserHomeDir")
	}
	certAbsolutePath := filepath.Join(homeDir, CertPath)
	keyAbsolutePath := filepath.Join(homeDir, KeyPath)
	if _, err = os.Stat(certAbsolutePath); os.IsNotExist(err) {
		return errors.Wrap(err, certNotExist.Error())
	}
	if _, err = os.Stat(keyAbsolutePath); os.IsNotExist(err) {
		return errors.Wrap(err, KeyNotExist.Error())
	}
	return nil
}
