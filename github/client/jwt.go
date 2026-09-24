package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	json "encoding/json/v2"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

// A runner clock running ahead of GitHub rejects a JWT issued in the future, and
// GitHub refuses one that lives longer than ten minutes.
const (
	JWTBackdate = 60 * time.Second
	JWTLifetime = 9 * time.Minute
)

type (
	jwtHeader struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}

	jwtClaims struct {
		IssuedAt  int64  `json:"iat"`
		ExpiresAt int64  `json:"exp"`
		Issuer    string `json:"iss"`
	}
)

// ParsePrivateKey reads an RSA key out of PEM, in the PKCS#1 form GitHub hands
// out as well as the PKCS#8 form a converted key ends up in.
func ParsePrivateKey(content []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(content)

	if block == nil {
		return nil, fmt.Errorf("Private key is not in PEM format.")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("Private key is neither a PKCS#1 nor a PKCS#8 key: %w", err)
	}

	key, ok := parsed.(*rsa.PrivateKey)

	if !ok {
		return nil, fmt.Errorf("Private key is not an RSA key.")
	}

	return key, nil
}

// SignJWT issues the RS256 JWT the app authenticates with to request an
// installation token.
func SignJWT(key *rsa.PrivateKey, id string, now time.Time) (string, error) {
	header, err := json.Marshal(jwtHeader{Algorithm: "RS256", Type: "JWT"})

	if err != nil {
		return "", err
	}

	claims, err := json.Marshal(jwtClaims{
		IssuedAt:  now.Add(-JWTBackdate).Unix(),
		ExpiresAt: now.Add(JWTLifetime).Unix(),
		Issuer:    id,
	})

	if err != nil {
		return "", err
	}

	unsigned := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(header),
		base64.RawURLEncoding.EncodeToString(claims),
	}, ".")

	digest := sha256.Sum256([]byte(unsigned))

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])

	if err != nil {
		return "", err
	}

	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}
