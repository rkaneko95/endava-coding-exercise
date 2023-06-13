package api

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"rkaneko/endava-coding-exercise/config"
	"time"
)

type Auth interface {
	CreateToken(authHeader string) (string, error)
}

type Service struct {
	Config        config.ServerConfig
	Log           *logrus.Logger
	TokenDuration time.Duration
	SecretKeyPath string
}

func (s *Service) CreateToken(authHeader string) (string, error) {
	claims, err := tokenClaims(authHeader, s.TokenDuration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	secretKey, err := readSecretKey(s.SecretKeyPath)
	if err != nil {
		return "", err
	}

	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}