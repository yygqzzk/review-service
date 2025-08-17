package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis 查询数据
func Get(ctx context.Context, rdb *redis.Client, key string) ([]byte, error) {
	return rdb.Get(ctx, key).Bytes()
}

// redis 缓存数据
func Set(ctx context.Context, rdb *redis.Client, key string, data []byte, expire time.Duration) error {
	return rdb.Set(ctx, key, data, expire).Err()
}
