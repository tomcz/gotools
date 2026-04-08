package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/huh/v2"
)

func main() {
	var certType string
	err := huh.NewSelect[string]().
		Title("Certificate type?").
		Options(huh.NewOption("Dev Root CA", "root"), huh.NewOption("TLS Certificate", "cert")).
		Value(&certType).
		Run()
	if err != nil {
		log.Fatalln(err)
	}
	if certType == "root" {
		err = createRootCA()
	} else {
		err = createCert()
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func createRootCA() error {
	commonName := "DevRootCA"
	country := "AU"
	organization := "Dev"
	locality := "Brisbane"
	province := "Queensland"
	streetAddress := "Adelaide Street"
	postalCode := "4000"
	keyType := "ec"
	certFilePath := "devca.crt"
	keyFilePath := "devca.key"
	certDaysTxt := "3650"

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Common Name?").Value(&commonName).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Country?").Value(&country).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Organization?").Value(&organization).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Locality?").Value(&locality).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Province?").Value(&province).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Street Address?").Value(&streetAddress).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Postal Code?").Value(&postalCode).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Certficate lifetime (in days)?").Value(&certDaysTxt).Validate(validPositiveInt),
			huh.NewSelect[string]().Title("Private key type?").Options(huh.NewOption("EC", "ec"), huh.NewOption("RSA", "rsa")).Value(&keyType),
			huh.NewInput().Title("Certificate file path?").Value(&certFilePath).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Private key file path?").Value(&keyFilePath).Validate(huh.ValidateNotEmpty()),
		),
	)
	err := form.Run()
	if err != nil {
		return err
	}
	certDays, err := strconv.Atoi(certDaysTxt)
	if err != nil {
		return err
	}

	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			CommonName:    commonName,
			Country:       []string{country},
			Organization:  []string{organization},
			Locality:      []string{locality},
			Province:      []string{province},
			StreetAddress: []string{streetAddress},
			PostalCode:    []string{postalCode},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, certDays),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	privKey, err := newPrivateKey(keyType)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKeyFor(keyType, privKey), privKey)
	if err != nil {
		return err
	}
	if err = writeCertificate(certBytes, certFilePath); err != nil {
		return err
	}
	if err = writePrivateKey(keyType, privKey, keyFilePath); err != nil {
		return err
	}
	return nil
}

func createCert() error {
	commonName := "DevCert"
	country := "AU"
	organization := "Dev"
	locality := "Brisbane"
	province := "Queensland"
	streetAddress := "Adelaide Street"
	postalCode := "4000"
	keyType := "ec"
	certFilePath := "dev.crt"
	keyFilePath := "dev.key"
	caCertFilePath := "devca.crt"
	caKeyFilePath := "devca.key"
	certDaysTxt := "365"
	var dnsNamesTxt string
	var ipAddrsTxt string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Common Name?").Value(&commonName).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Country?").Value(&country).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Organization?").Value(&organization).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Locality?").Value(&locality).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Province?").Value(&province).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Street Address?").Value(&streetAddress).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Postal Code?").Value(&postalCode).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Certficate lifetime (in days):").Value(&certDaysTxt).Validate(validPositiveInt),
			huh.NewSelect[string]().Title("Private key type?").Options(huh.NewOption("EC", "ec"), huh.NewOption("RSA", "rsa")).Value(&keyType),
			huh.NewInput().Title("DNS SANs (comma-separated):").Value(&dnsNamesTxt),
			huh.NewInput().Title("IP SANs (comma-separated):").Value(&ipAddrsTxt),
			huh.NewInput().Title("Certificate file path?").Value(&certFilePath).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Private key file path?").Value(&keyFilePath).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Root certificate file path?").Value(&caCertFilePath).Validate(huh.ValidateNotEmpty()),
			huh.NewInput().Title("Root key file path?").Value(&caKeyFilePath).Validate(huh.ValidateNotEmpty()),
		),
	)
	err := form.Run()
	if err != nil {
		return err
	}
	certDays, err := strconv.Atoi(certDaysTxt)
	if err != nil {
		return err
	}
	var dnsNames []string
	if dnsNamesTxt != "" {
		for _, name := range strings.Split(dnsNamesTxt, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				dnsNames = append(dnsNames, name)
			}
		}
	}
	var ipAddrs []net.IP
	if ipAddrsTxt != "" {
		for _, value := range strings.Split(ipAddrsTxt, ",") {
			value = strings.TrimSpace(value)
			if value != "" {
				ip := net.ParseIP(value)
				if ip == nil {
					return fmt.Errorf("invalid IP address: %s", value)
				}
				ipAddrs = append(ipAddrs, ip)
			}
		}
	}
	caCert, err := readCertificate(caCertFilePath)
	if err != nil {
		return err
	}
	caKey, err := readPrivateKey(caKeyFilePath)
	if err != nil {
		return err
	}

	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			CommonName:    commonName,
			Country:       []string{country},
			Organization:  []string{organization},
			Locality:      []string{locality},
			Province:      []string{province},
			StreetAddress: []string{streetAddress},
			PostalCode:    []string{postalCode},
		},
		DNSNames:    dnsNames,
		IPAddresses: ipAddrs,
		NotBefore:   now,
		NotAfter:    now.AddDate(0, 0, certDays),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	privKey, err := newPrivateKey(keyType)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, publicKeyFor(keyType, privKey), caKey)
	if err != nil {
		return err
	}
	if err = writeCertificate(certBytes, certFilePath); err != nil {
		return err
	}
	if err = writePrivateKey(keyType, privKey, keyFilePath); err != nil {
		return err
	}
	return nil
}

func validPositiveInt(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num < 1 {
		return errors.New("minimum value is 1")
	}
	return nil
}

func newPrivateKey(keyType string) (any, error) {
	if keyType == "ec" {
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		return priv, err
	}
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	return priv, err
}

func publicKeyFor(keyType string, privKey any) any {
	if keyType == "ec" {
		pk := privKey.(*ecdsa.PrivateKey)
		return &pk.PublicKey
	}
	pk := privKey.(*rsa.PrivateKey)
	return &pk.PublicKey
}

func writeCertificate(certBytes []byte, outfile string) error {
	log.Println("Writing to", outfile)
	fp, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer fp.Close()

	return pem.Encode(fp, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
}

func readCertificate(certFile string) (*x509.Certificate, error) {
	buf, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	for {
		var block *pem.Block
		block, buf = pem.Decode(buf)
		if block == nil {
			return nil, fmt.Errorf("CERTIFICATE block not found in %s", certFile)
		}
		if block.Type == "CERTIFICATE" {
			return x509.ParseCertificate(block.Bytes)
		}
	}
}

func writePrivateKey(keyType string, privKey any, outfile string) error {
	log.Println("Writing to", outfile)
	fp, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer fp.Close()

	if keyType != "ec" {
		pk := privKey.(*rsa.PrivateKey)
		return pem.Encode(fp, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(pk),
		})
	}

	pk := privKey.(*ecdsa.PrivateKey)
	buf, err := x509.MarshalECPrivateKey(pk)
	if err != nil {
		return err
	}
	return pem.Encode(fp, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: buf,
	})
}

func readPrivateKey(keyFile string) (any, error) {
	buf, fail := os.ReadFile(keyFile)
	if fail != nil {
		return nil, fail
	}
	for {
		var block *pem.Block
		block, buf = pem.Decode(buf)
		if block == nil {
			return nil, fmt.Errorf("(RSA|EC) PRIVATE KEY block not found in %s", keyFile)
		}
		if block.Type == "RSA PRIVATE KEY" {
			return x509.ParsePKCS1PrivateKey(block.Bytes)
		}
		if block.Type == "EC PRIVATE KEY" {
			return x509.ParseECPrivateKey(block.Bytes)
		}
	}
}
