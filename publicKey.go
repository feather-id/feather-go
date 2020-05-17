package feather

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
)

func (s *sessions) getPublicKey(keyID string) (*rsa.PublicKey, error) {

	// Check the cache
	if publicKey, ok := s.cachedPublicKeys[keyID]; ok {
		return publicKey, nil
	}

	// Query Feather API for the key
	type publicKeyResponse struct {
		ID     string `json:"id"`
		Object string `json:"object"`
		PEM    string `json:"pem"`
	}
	var pubKeyResponse publicKeyResponse
	path := strings.Join([]string{pathPublicKeys, keyID}, "/")
	if err := s.gateway.sendRequest(http.MethodGet, path, nil, &pubKeyResponse); err != nil {
		return nil, err
	}

	// Decode and parse the key
	pubPem, _ := pem.Decode([]byte(pubKeyResponse.PEM))
	if pubPem == nil {
		return nil, fmt.Errorf("Failed to parse public key %v", keyID)
	}
	if pubPem.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Decoded key is of the wrong type (%v)", pubPem.Type)
	}
	var parsedKey interface{}
	var err error
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS1PublicKey(pubPem.Bytes); err != nil {
			return nil, err
		}
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to parse public key %v", keyID)
	}

	// Cache and return
	s.cachedPublicKeys[keyID] = publicKey
	return publicKey, nil
}
