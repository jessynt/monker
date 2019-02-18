package storage

import (
	"net/url"

	"github.com/go-redis/redis"
)

type RedisStorage struct {
	redisClient *redis.Client
	key         string
}

func NewRedisStorage(dsn *url.URL, key string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: dsn.Host,
	})
	return &RedisStorage{redisClient: client, key: key}, nil
}

func (s *RedisStorage) Put(data []byte) error {
	return s.redisClient.RPush(s.key, data).Err()
}

func (s *RedisStorage) Get() ([]byte, error) {
	bytes, err := s.redisClient.LPop(s.key).Bytes()
	if err == redis.Nil {
		return nil, ErrStorageIsEmpty
	}
	return bytes, err
}

func (s *RedisStorage) Close() error {
	return s.redisClient.Close()
}
