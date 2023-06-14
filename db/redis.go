package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Client *redis.Client
}

type RedisClient interface {
	GetBytes(key string) ([]byte, error)
	SetBytes(key string, content []byte) error
	GetSignatureKeys() ([]string, error)
}

func RunRedis(password, host string, port int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r RedisService) GetBytes(key string) ([]byte, error) {
	value, err := r.Client.Get(r.Client.Context(), key).Bytes()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (r RedisService) SetBytes(key string, content []byte) error {
	err := r.Client.Set(r.Client.Context(), key, content, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisService) GetSignatureKeys() ([]string, error) {
	value, err := r.Client.Keys(r.Client.Context(), "signature_*").Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}
