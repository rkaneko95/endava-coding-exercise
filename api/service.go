package api

import (
	"errors"
	"fmt"
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

func (s *Service) VerifyToken(authHeader string) (*Claims, error) {
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
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, errors.New("token err: token is invalid")
	}

	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
