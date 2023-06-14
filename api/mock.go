package api

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"os"
	"rkaneko/endava-coding-exercise/db"
	"strings"
)

func MockedGetPrivateKey(key string, rds *db.RedisService) (*rsa.PrivateKey, error) {
	fileName := strings.TrimPrefix(key, "secret_")
	file, err := os.Open(fileName)
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

func MockedGetPublicKey(key string, rds *db.RedisService) (*rsa.PublicKey, error) {
	fileName := strings.TrimPrefix(key, "signature_")
	file, err := os.Open(fmt.Sprintf("%s.pub", fileName))
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
