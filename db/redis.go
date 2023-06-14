package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Client *redis.Client
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

func (r RedisService) GetString(key string) (string, error) {
	value, err := r.Client.Get(r.Client.Context(), key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r RedisService) GetSignatureKeys() ([]string, error) {
	value, err := r.Client.Keys(r.Client.Context(), "signature_*").Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}
