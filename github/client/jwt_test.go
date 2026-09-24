package client

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	json "encoding/json/v2"
	"encoding/pem"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/github/test/keys"
)

var _ = Describe("ParsePrivateKey", func() {
	It("reads the PKCS#1 key GitHub hands out", func() {
		key, content := keys.Generate()

		parsed, err := ParsePrivateKey([]byte(content))
		Expect(err).NotTo(HaveOccurred())

		Expect(parsed.Equal(key)).To(BeTrue())
	})

	It("reads a key converted to PKCS#8", func() {
		key, _ := keys.Generate()

		der, err := x509.MarshalPKCS8PrivateKey(key)
		Expect(err).NotTo(HaveOccurred())

		// a key parsed out of PKCS#8 carries the same numbers in a different internal
		// layout, which only the key's own Equal sees through.
		parsed, err := ParsePrivateKey(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
		Expect(err).NotTo(HaveOccurred())

		Expect(parsed.Equal(key)).To(BeTrue())
	})

	It("rejects contents that are not PEM", func() {
		_, err := ParsePrivateKey([]byte("not a key"))

		Expect(err).To(MatchError("Private key is not in PEM format."))
	})

	It("rejects a PEM block that is not a private key", func() {
		_, err := ParsePrivateKey(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: []byte("garbage")}))

		Expect(err).To(MatchError(ContainSubstring("neither a PKCS#1 nor a PKCS#8 key")))
	})

	// GitHub only signs with RS256, so a key of any other kind can never mint a token.
	It("rejects a key that is not RSA", func() {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		Expect(err).NotTo(HaveOccurred())

		der, err := x509.MarshalPKCS8PrivateKey(key)
		Expect(err).NotTo(HaveOccurred())

		_, err = ParsePrivateKey(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))

		Expect(err).To(MatchError("Private key is not an RSA key."))
	})
})

var _ = Describe("SignJWT", func() {
	var (
		key *rsa.PrivateKey
		now = time.Unix(1_800_000_000, 0)
	)

	BeforeEach(func() {
		key, _ = keys.Generate()
	})

	decode := func(segment string, v any) {
		GinkgoHelper()

		raw, err := base64.RawURLEncoding.DecodeString(segment)
		Expect(err).NotTo(HaveOccurred())
		Expect(json.Unmarshal(raw, v)).To(Succeed())
	}

	It("signs with RS256 so the public half of the key verifies it", func() {
		jwt, err := SignJWT(key, "123", now)
		Expect(err).NotTo(HaveOccurred())

		segments := strings.Split(jwt, ".")
		Expect(segments).To(HaveLen(3))

		signature, err := base64.RawURLEncoding.DecodeString(segments[2])
		Expect(err).NotTo(HaveOccurred())

		digest := sha256.Sum256([]byte(segments[0] + "." + segments[1]))
		Expect(rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, digest[:], signature)).To(Succeed())
	})

	It("carries the header GitHub expects", func() {
		jwt, err := SignJWT(key, "123", now)
		Expect(err).NotTo(HaveOccurred())

		header := map[string]string{}
		decode(strings.Split(jwt, ".")[0], &header)

		Expect(header).To(Equal(map[string]string{"alg": "RS256", "typ": "JWT"}))
	})

	// the issue time is backdated against a runner clock ahead of GitHub, and the
	// whole window stays under the ten minutes GitHub accepts.
	It("issues the claims for the app inside the window GitHub accepts", func() {
		jwt, err := SignJWT(key, "123", now)
		Expect(err).NotTo(HaveOccurred())

		claims := map[string]any{}
		decode(strings.Split(jwt, ".")[1], &claims)

		Expect(claims).To(Equal(map[string]any{
			"iat": float64(now.Unix() - 60),
			"exp": float64(now.Unix() + 540),
			"iss": "123",
		}))
	})

	It("leaves the padding off every segment", func() {
		jwt, err := SignJWT(key, "123", now)
		Expect(err).NotTo(HaveOccurred())

		Expect(jwt).NotTo(ContainSubstring("="))
	})
})
