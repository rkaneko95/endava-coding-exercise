package api

import (
	"encoding/base64"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTokenClaims(t *testing.T) {
	l := gofakeit.Number(25, 50)
	randBytes := make([]byte, l)
	base64Str := base64.StdEncoding.EncodeToString(randBytes)

	duration := time.Hour

	expectedIssuedAt := time.Now()
	expectedExpiredAt := expectedIssuedAt.Add(duration)

	c, err := tokenClaims(fmt.Sprintf("Basic %s", base64Str), duration)

	assert.NoError(t, err)
	assert.NotEmpty(t, c)
	assert.IsType(t, &Claims{}, c)
	assert.Equal(t, base64Str, c.Base64)
	assert.WithinDuration(t, expectedExpiredAt, c.ExpiredAt, time.Second)
	assert.WithinDuration(t, expectedIssuedAt, c.IssuedAt, time.Second)
}

func TestTokenClaimsInvalidAuth(t *testing.T) {
	c, err := tokenClaims("error", time.Hour)

	assert.Error(t, err)
	assert.Empty(t, c)
}

func TestValid(t *testing.T) {
	l := gofakeit.Number(25, 50)
	randBytes := make([]byte, l)
	base64Str := base64.StdEncoding.EncodeToString(randBytes)

	authHeader := fmt.Sprintf("Basic %s", base64Str)

	duration := time.Hour

	c, _ := tokenClaims(authHeader, duration)

	err := c.Valid()
	assert.NoError(t, err)
}

func TestValidTokenExpired(t *testing.T) {
	l := gofakeit.Number(25, 50)
	randBytes := make([]byte, l)
	base64Str := base64.StdEncoding.EncodeToString(randBytes)

	authHeader := fmt.Sprintf("Basic %s", base64Str)

	duration := time.Hour * -1

	c, _ := tokenClaims(authHeader, duration)

	err := c.Valid()
	assert.Error(t, err)
}
