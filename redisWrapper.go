package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Rdb struct {
	rdb *redis.Client
}
type Config struct {
	Addr         string
	Username     string
	Password     string
	DB           int
	MaxRetries   int
	DialTimeout  time.Duration
	PoolSize     int
	MinIdleConns int
}

//go:generate mockgen --destination=./mocks/mock_Cacher.go --package=mocks github.com/vatsal278/go-redis-cache Cacher
type Cacher interface {
	Get(string) ([]byte, error)
	Set(string, interface{}, time.Duration) error
	Health() (string, error)
}

func NewCacher(c Config) Cacher {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Username:     c.Username,
		Password:     c.Password, // no password set
		DB:           c.DB,       // use default DB
		MaxRetries:   c.MaxRetries,
		DialTimeout:  c.DialTimeout,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
	})
	return &Rdb{
		rdb: rdb,
	}
}
func (r Rdb) Get(key string) ([]byte, error) {
	data, err := r.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	return data, err
}
func (r Rdb) Set(key string, value interface{}, expiry time.Duration) error {
	err := r.rdb.Set(context.Background(), key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r Rdb) Health() (string, error) {
	status := r.rdb.Ping(context.Background())
	return status.Result()
}
