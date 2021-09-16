package main

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func run() error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	for block, data := pem.Decode(input); block != nil; block, data = pem.Decode(data) {
		dumpPEMBlock(block)
	}

	return nil
}

func dumpPEMBlock(block *pem.Block) {
	var (
		value interface{}
		err   error
	)
	fmt.Fprintln(os.Stderr, block.Type)

	switch bytes := block.Bytes; block.Type {
	case "CERTIFICATE":
		value, err = x509.ParseCertificate(bytes)
	case "CERTIFICATE REQUEST":
		value, err = x509.ParseCertificateRequest(bytes)
	case "PRIVATE KEY":
		value, err = x509.ParsePKCS8PrivateKey(bytes)
	case "PUBLIC KEY":
		value, err = x509.ParsePKIXPublicKey(bytes)
	case "RSA PRIVATE KEY":
		value, err = x509.ParsePKCS1PrivateKey(bytes)
	case "EC PRIVATE KEY":
		value, err = x509.ParseECPrivateKey(bytes)
	case "RSA PUBLIC KEY":
		value, err = x509.ParsePKCS1PublicKey(bytes)
	case "X509 CRL":
		value, err = x509.ParseDERCRL(bytes)
	default:
		var rest []byte
		rest, err = asn1.Unmarshal(bytes, &value)
		if len(rest) > 0 {
			result := []interface{}{value}
			var val interface{}
			for len(rest) > 0 && err == nil {
				rest, err = asn1.Unmarshal(rest, &val)
				if val != nil {
					result = append(result, val)
					val = nil
				}
			}
			value = result
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	if value != nil {
		spew.Fdump(os.Stderr, value)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
