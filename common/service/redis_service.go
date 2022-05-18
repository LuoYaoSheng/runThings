package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisKeyExpiredFunc func(string)

var ctx = context.Background()
var redisClient *redis.Client

func GetRedisClient(addr, password string, db int) *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return redisClient
}

func SetRdValue(key string, value interface{}) error {
	if redisClient == nil {
		return errors.New("redis客户端连接失败")
	}
	err := redisClient.Set(ctx, key, value, 0).Err()
	return err
}

func SetRdValueTimeout(key string, value interface{}, expiration time.Duration) error {
	if redisClient == nil {
		return errors.New("redis客户端连接失败")
	}
	err := redisClient.Set(ctx, key, value, expiration).Err()
	return err
}

func GetRdValue(key string) (string, error) {
	if redisClient == nil {
		return "", errors.New("redis客户端连接失败")
	}
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", err
	} else {
		return val, err
	}
}

func DelRdValue(key string) (int64, error) {
	if redisClient == nil {
		return 0, errors.New("redis客户端连接失败")
	}
	return redisClient.Del(ctx, key).Result()
}

func SubscribeKeyExpired(fc RedisKeyExpiredFunc) error {
	if redisClient == nil {
		return errors.New("redis客户端连接失败")
	}
	go func() {
		pubsub := redisClient.Subscribe(ctx, "__keyevent@0__:expired")
		defer pubsub.Close()
		for msg := range pubsub.Channel() {
			fc(msg.Payload)
		}
	}()
	return nil
}
