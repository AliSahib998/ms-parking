package cache

import (
	"context"
	"fmt"
	"github.com/AliSahib998/ms-parking/config"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

const connectionPoolSize = 4

type IRedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd
}

// NewRedisClusterClient for prod env
func NewRedisClusterClient() *redis.ClusterClient {
	var addrs []string
	for _, v := range config.Props.RedisURL {
		v = strings.TrimSpace(v)
		o, _ := redis.ParseURL(v)
		addrs = append(addrs, o.Addr)
	}

	fmt.Print(addrs)

	opts := &redis.ClusterOptions{
		Addrs:        addrs,
		MaxRedirects: 0,
		Password:     "",
		PoolSize:     connectionPoolSize,
	}

	return redis.NewClusterClient(opts)
}

func NewRedisClientForDev() *redis.Client {
	var addrs []string
	for _, v := range config.Props.RedisURL {
		v = strings.TrimSpace(v)
		o, _ := redis.ParseURL(v)
		addrs = append(addrs, o.Addr)
	}

	fmt.Print(addrs)

	opts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	return redis.NewClient(opts)
}
