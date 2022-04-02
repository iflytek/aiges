package storage

import (
	"github.com/go-redis/redis"
	"time"
)

//redis
type RedisConfig struct {
	Addrs        []string
	Password     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	DialTimeout  time.Duration
	PoolSize     int
}

func NewRedis(conf RedisConfig) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        conf.Addrs,
		Password:     conf.Password,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
		DialTimeout:  conf.DialTimeout,
		PoolSize:     conf.PoolSize,
	})
}
