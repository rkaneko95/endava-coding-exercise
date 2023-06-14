package api

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Base64    string    `json:"base64"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func tokenClaims(authHeader string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	credBase64, err := extractAuth(authHeader)
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Base64:    credBase64,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (c Payload) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return errors.New("token err: token has expired")
	}
	return nil
}
