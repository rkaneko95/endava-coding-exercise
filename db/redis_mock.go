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
	privateKey, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
