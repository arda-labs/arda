package redis

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(addr string, db int, password string, logger log.Logger) (*RedisClient, func(), error) {
	l := log.NewHelper(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}

	l.Info("redis client initialized")

	cleanup := func() {
		rdb.Close()
		l.Info("redis client closed")
	}

	return &RedisClient{Client: rdb}, cleanup, nil
}
