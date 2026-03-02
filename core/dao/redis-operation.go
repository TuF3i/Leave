package dao

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *Dao) checkKeyExistence(ctx context.Context, key string) (bool, error) {
	_, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Dao) getKeyValue(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *Dao) setNewValue(ctx context.Context, key string, value string, expireTime time.Duration) error {
	err := r.rdb.Set(ctx, key, value, expireTime).Err()
	return err
}
