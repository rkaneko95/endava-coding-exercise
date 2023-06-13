package api

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	ID        uuid.UUID `json:"id"`
	Base64    string    `json:"base64"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func tokenClaims(authHeader string, duration time.Duration) (*Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	credBase64, err := extractAuth(authHeader)
	if err != nil {
		return nil, err
	}

	claims := &Claims{
		ID:        tokenID,
		Base64:    credBase64,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return claims, nil
}

func (c Claims) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return errors.New("token err: token has expired")
	}
	return nil
}
