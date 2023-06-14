package api

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	monkey.Patch(getPrivateKey, MockedGetPrivateKey)
	monkey.Patch(getPublicKey, MockedGetPublicKey)
}

func TestCreateToken(t *testing.T) {
	cases := []struct {
		service *Service
		input   string
		hasErr  bool
	}{
		{
			service: &Service{
				TokenDuration: 24 * time.Hour,
				KeyUUID:       "../test_resources/mock_RS256.key",
			},
			input:  "Basic randumstring==",
			hasErr: false,
		},
		{
			service: &Service{
				TokenDuration: 24 * time.Hour,
				KeyUUID:       "",
			},
			input:  "Basic randumstring==",
			hasErr: true,
		},
		{
			service: &Service{
				TokenDuration: 24 * time.Hour,
				KeyUUID:       "../test_resources/mock_RS256.key",
			},
			input:  "randumstring==",
			hasErr: true,
		},
	}

	for _, c := range cases {
		token, err := c.service.CreateToken(c.input)
		if c.hasErr {
			assert.Error(t, err)
			assert.Empty(t, token)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		}
	}
}

func TestVerifyToken(t *testing.T) {
	cases := []struct {
		service *Service
		hasErr  bool
	}{
		{
			service: &Service{
				TokenDuration: 24 * time.Hour,
				KeyUUID:       "../test_resources/mock_RS256.key",
			},
			hasErr: false,
		},
		{
			service: &Service{
				TokenDuration: 24 * time.Hour * -1,
				KeyUUID:       "../test_resources/mock_RS256.key",
			},
			hasErr: true,
		},
	}

	for _, c := range cases {
		input, _ := c.service.CreateToken("Basic randumstring==")

		issue, expired, err := c.service.VerifyToken(fmt.Sprintf("Bearer %s", input))
		if c.hasErr {
			assert.Error(t, err)
			assert.Nil(t, issue)
			assert.Nil(t, expired)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, issue)
			assert.NotEmpty(t, expired)
		}
	}
}
