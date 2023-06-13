package api

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"os"
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

func getPrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	keyData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKey(path string) (*rsa.PublicKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	keyData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
