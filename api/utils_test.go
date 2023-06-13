package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractAuth(t *testing.T) {
	cases := []struct {
		input       string
		expectedStr string
		hasErr      bool
	}{
		{
			input:       "Basic randumstring==",
			expectedStr: "randumstring==",
			hasErr:      false,
		},
		{
			input:       "randumstring==",
			expectedStr: "",
			hasErr:      true,
		},
		{
			input:       "",
			expectedStr: "",
			hasErr:      true,
		},
	}

	for _, c := range cases {
		result, err := extractAuth(c.input)
		assert.Equal(t, c.expectedStr, result)
		if c.hasErr {
			assert.NotEmpty(t, err)
		} else {
			assert.Empty(t, err)
		}
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	cases := []struct {
		input       string
		expectedStr string
	}{
		{
			input:       "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedStr: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			input:       "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedStr: "",
		},
		{
			input:       "",
			expectedStr: "",
		},
	}

	for _, c := range cases {
		result := extractTokenFromHeader(c.input)
		assert.Equal(t, c.expectedStr, result)
	}
}

func TestGetPrivateKey(t *testing.T) {
	cases := []struct {
		input  string
		hasErr bool
	}{
		{
			input:  "../test_resources/mock_RS256.key",
			hasErr: false,
		},
		{
			input:  "../test_resources/no_mock.key",
			hasErr: true,
		},
		{
			input:  "",
			hasErr: true,
		},
	}

	for _, c := range cases {
		p, err := getPrivateKey(c.input)

		if c.hasErr {
			assert.Error(t, err)
			assert.Nil(t, p)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, p)
		}
	}
}

func TestGetPublicKey(t *testing.T) {
	cases := []struct {
		input  string
		hasErr bool
	}{
		{
			input:  "../test_resources/mock_RS256.key.pub",
			hasErr: false,
		},
		{
			input:  "../test_resources/no_mock.key.pub",
			hasErr: true,
		},
		{
			input:  "",
			hasErr: true,
		},
	}

	for _, c := range cases {
		p, err := getPublicKey(c.input)

		if c.hasErr {
			assert.Error(t, err)
			assert.Nil(t, p)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, p)
		}
	}
}
