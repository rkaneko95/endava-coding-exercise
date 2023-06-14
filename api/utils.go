package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/square/go-jose/v3"
	"rkaneko/endava-coding-exercise/db"
	"strings"
)

func extractAuth(header string) (string, error) {
	if !strings.HasPrefix(header, "Basic ") {
		return "", errors.New("invalid Authorization header format")
	}

	credBase64 := strings.TrimPrefix(header, "Basic ")
	return credBase64, nil
}

func extractTokenFromHeader(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func getPrivateKey(key string, rds *db.RedisService) (*rsa.PrivateKey, error) {
	data, err := rds.GetBytes(key)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: data,
	})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKey(key string, rds *db.RedisService) (*rsa.PublicKey, error) {
	data, err := rds.GetBytes(key)
	if err != nil {
		return nil, err
	}

	var jwk jose.JSONWebKey
	err = jwk.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(jwk.Key)
	if err != nil {
		return nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
