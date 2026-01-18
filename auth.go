package rasedi

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type Auth struct {
	privateKeyPEM string
	keyID         string
}

func NewAuth(privateKeyPEM, keyID string) *Auth {
	return &Auth{
		privateKeyPEM: privateKeyPEM,
		keyID:         keyID,
	}
}

func (a *Auth) GetKeyID() string {
	return a.keyID
}

func (a *Auth) MakeSignature(method, relativeURL string) (string, error) {
	rawSign := fmt.Sprintf("%s || %s || %s", method, a.keyID, relativeURL)
	data := []byte(rawSign)

	block, _ := pem.Decode([]byte(a.privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the key")
	}

	var privateKey interface{}
	var err error

	if privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			if privateKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
				return "", fmt.Errorf("failed to parse private key: %v", err)
			}
		}
	}

	var signature []byte

	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		hashed := sha256.Sum256(data)
		signature, err = rsa.SignPKCS1v15(rand.Reader, k, crypto.SHA256, hashed[:])
		if err != nil {
			return "", fmt.Errorf("rsa sign failed: %v", err)
		}
	case *ecdsa.PrivateKey:
		hashed := sha256.Sum256(data)
		signature, err = ecdsa.SignASN1(rand.Reader, k, hashed[:])
		if err != nil {
			return "", fmt.Errorf("ecdsa sign failed: %v", err)
		}
	case ed25519.PrivateKey:
		signature = ed25519.Sign(k, data)
	default:
		return "", fmt.Errorf("unsupported private key type: %T", k)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
