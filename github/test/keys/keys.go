// Package keys hands the specs a GitHub App private key of their own, so nothing
// real ever has to sit in the repository.
package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Generate renders a fresh RSA key as PKCS#1 PEM, the form GitHub hands out.
func Generate() (*rsa.PrivateKey, string) {
	GinkgoHelper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())

	return key, string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}))
}
