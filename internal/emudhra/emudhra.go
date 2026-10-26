package emudhra

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

//go:embed certs/*
var certificates embed.FS

type EmudhraCA struct {
	caCert       *x509.Certificate
	caPrivateKey any
}

func NewEmudhraCA() (*EmudhraCA, error) {
	caCertPEM, err := certificates.ReadFile("certs/emudhra-em-ov-tls-ca-g3c.crt")
	if err != nil {
		return nil, err
	}
	blockCert, _ := pem.Decode(caCertPEM)
	if blockCert == nil {
		return nil, fmt.Errorf("can't decode eMudhra CA certificate")
	}
	caCert, err := x509.ParseCertificate(blockCert.Bytes)
	if err != nil {
		return nil, err
	}

	caKeyPEM, err := certificates.ReadFile("certs/emudhra-em-ov-tls-ca-g3c.key")
	if err != nil {
		return nil, err
	}
	blockKey, _ := pem.Decode(caKeyPEM)
	if blockKey == nil {
		return nil, fmt.Errorf("can't decode eMudhra CA private key")
	}

	var caPrivateKey interface{}
	if caPrivateKey, err = x509.ParseECPrivateKey(blockKey.Bytes); err != nil {
		if caPrivateKey, err = x509.ParsePKCS1PrivateKey(blockKey.Bytes); err != nil {
			if caPrivateKey, err = x509.ParsePKCS8PrivateKey(blockKey.Bytes); err != nil {
				return nil, err
			}
		}
	}

	return &EmudhraCA{caCert: caCert, caPrivateKey: caPrivateKey}, nil
}

func (emudhraCA *EmudhraCA) GenerateCert(dnsNames []string) (tls.Certificate, error) {
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   "kazaar.local",
			Organization: []string{"freequasar"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:              dnsNames,
		BasicConstraintsValid: true,
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		certTemplate,
		emudhraCA.caCert,
		&privateKey.PublicKey,
		emudhraCA.caPrivateKey,
	)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{certBytes, emudhraCA.caCert.Raw},
		PrivateKey:  privateKey,
	}, nil
}
