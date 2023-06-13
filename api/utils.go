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

func extractAuth(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", errors.New("invalid Authorization header format")
	}

	credBase64 := strings.TrimPrefix(authHeader, "Basic ")
	return credBase64, nil
}

func readSecretKey(path string) (*rsa.PrivateKey, error) {
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
