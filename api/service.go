package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/square/go-jose/v3"
	"rkaneko/endava-coding-exercise/config"
	"rkaneko/endava-coding-exercise/db"
	"time"
)

type Auth interface {
	CreateToken(authHeader string) (string, error)
	VerifyToken(authHeader string) (*time.Time, *time.Time, error)
	ListSigningKeys() ([]jose.JSONWebKey, error)
}

type Service struct {
	Config        config.ServerConfig
	Log           *logrus.Logger
	TokenDuration time.Duration
	KeyUUID       string
	RedisService  *db.RedisService
}

func (s *Service) CreateToken(authHeader string) (string, error) {
	claims, err := tokenClaims(authHeader, s.TokenDuration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	secretKey, err := getPrivateKey(
		fmt.Sprintf("secret_%s", s.KeyUUID),
		s.RedisService)
	if err != nil {
		return "", err
	}

	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) VerifyToken(authHeader string) (*time.Time, *time.Time, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		secretKey, err := getPublicKey(
			fmt.Sprintf("signature_%s", s.KeyUUID),
			s.RedisService)
		if err != nil {
			return "", err
		}
		return secretKey, nil
	}

	token := extractTokenFromHeader(authHeader)

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, nil, errors.New("token err: token is invalid")
	}

	return &claims.IssuedAt, &claims.ExpiredAt, nil
}

func (s *Service) ListSigningKeys() ([]jose.JSONWebKey, error) {
	keys, err := s.RedisService.GetSignatureKeys()
	if err != nil {
		return nil, err
	}

	signatures := make([]jose.JSONWebKey, 0, len(keys))
	for _, key := range keys {
		data, err := s.RedisService.GetBytes(key)
		if err != nil {
			return nil, err
		}

		var signature jose.JSONWebKey
		err = json.Unmarshal(data, &signature)
		if err != nil {
			return nil, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}
