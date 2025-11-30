package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

type RedisStoreConfig struct {
	Address   string
	Passwd    string
	DefaultDb int
}

// NewRedisStore initializes the connection
func NewRedisStore(redisConfig RedisStoreConfig) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,   // e.g., "localhost:6379"
		Password: redisConfig.Passwd,    // no password set
		DB:       redisConfig.DefaultDb, // use default DB
	})

	return &RedisStore{client: client}
}

func (r *RedisStore) GetObject(key string) (string, error) {
	keyVal := r.client.Get(context.Background(), key)
	if keyVal.Err() != nil {
		return "", fmt.Errorf("object not found in the db")
	}
	return keyVal.Result()
}

func (r *RedisStore) SetObject(key string, value string) error {
	if keyVal := r.client.Get(context.Background(), key); keyVal != nil {
		fmt.Println("rewriting the key object")
	}
	if err := r.client.Set(context.Background(), key, string(value), 0).Err(); err != nil {
		return fmt.Errorf("faield to store object in to keystore %s", err)
	}
	return nil
}
