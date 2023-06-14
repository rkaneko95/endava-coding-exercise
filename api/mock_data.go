package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/square/go-jose/v3"
)

func (s *Service) GenerateKeys() error {
	id, err := s.generateUUID()
	if err != nil {
		return err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	if err = s.setPrivateKey(privateKey, id); err != nil {
		return err
	}

	if err = s.setPublicKey(privateKey, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) generateUUID() (string, error) {
	keys, err := s.RedisService.GetSignatureKeys()
	if err != nil || len(keys) == 0 {
		return s.KeyUUID, nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *Service) setPrivateKey(privateKey *rsa.PrivateKey, id string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	err := s.RedisService.SetBytes(
		fmt.Sprintf("secret_%s", id),
		privateKeyBytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) setPublicKey(privateKey *rsa.PrivateKey, id string) error {
	jwk := jose.JSONWebKey{Key: &privateKey.PublicKey}
	jsonBytes, err := json.Marshal(jwk)
	if err != nil {
		return err
	}

	err = s.RedisService.SetBytes(
		fmt.Sprintf("signature_%s", id),
		jsonBytes)
	if err != nil {
		return err
	}

	return nil
}
