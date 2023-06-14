package db

import (
	"bufio"
	"io"
	"os"
)

func (r RedisService) MockedGetBytes(key string) ([]byte, error) {
	file, err := os.Open(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	pKey, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return pKey, nil
}

func (r RedisService) MockedGetSignatureKeys() ([]string, error) {
	return []string{"signature_1", "signature_2"}, nil
}
