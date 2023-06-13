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
