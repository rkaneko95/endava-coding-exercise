package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"rkaneko/endava-coding-exercise/config"
	"rkaneko/endava-coding-exercise/db"
	"time"
)

type Auth interface {
	CreateToken(authHeader string) (string, error)
	VerifyToken(authHeader string) (*time.Time, *time.Time, error)
	ListSigningKeys() ([]Signature, error)
}

type Service struct {
	Config        config.ServerConfig
	Log           *logrus.Logger
	TokenDuration time.Duration
	SecretKeyPath string
	RedisService  *db.RedisService
}

func (s *Service) CreateToken(authHeader string) (string, error) {
	claims, err := tokenClaims(authHeader, s.TokenDuration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	secretKey, err := getPrivateKey(s.SecretKeyPath)
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
		secretKey, err := getPublicKey(fmt.Sprintf("%s.pub", s.SecretKeyPath))
		if err != nil {
			return "", err
		}
		return secretKey, nil
	}

	token := extractTokenFromHeader(authHeader)

	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, nil, errors.New("token err: token is invalid")
	}

	return &claims.IssuedAt, &claims.ExpiredAt, nil
}

func (s *Service) ListSigningKeys() ([]Signature, error) {
	keys, err := s.RedisService.GetSignatureKeys()
	if err != nil {
		return nil, err
	}

	signatures := make([]Signature, 0, len(keys))
	for _, key := range keys {
		data, err := s.RedisService.GetString(key)
		if err != nil {
			return nil, err
		}

		var signature Signature
		err = json.Unmarshal([]byte(data), &signature)
		if err != nil {
			return nil, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}
