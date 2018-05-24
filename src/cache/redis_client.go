package cache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	REDIS_ADDRESS = "127.0.0.1"
	REDIS_PORT    = 6379
	TOKEN_PREFIX  = "token"
)

var client redis.Conn

// InitRedisClient 初始化redis连接
func InitRedisClient() (err error) {
	client, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", REDIS_ADDRESS, REDIS_PORT))
	return
}

// TokenCache 缓存token并设置相应的ttl
func TokenCache(token string, account string, ttl int64) (err error) {
	key := fmt.Sprintf("%s_%s", TOKEN_PREFIX, token)
	_, err = client.Do("SET", key, account, "EX", ttl)
	return
}

// GetAndRefreshToken 获取token对应的账号并更新ttl,如果不存在记录则返回nil
func GetAndRefreshToken(token string, ttl int64) (account *string, err error) {
	key := fmt.Sprintf("%s_%s", TOKEN_PREFIX, token)
	fmt.Println("key = ", key)
	if exists, err := redis.Bool(client.Do("EXISTS", key)); err == nil && exists {
		if _, err = client.Do("EXPIRE", key, ttl); err == nil {
			if value, err := redis.String(client.Do("GET", key)); err == nil {
				return &value, nil
			}
			fmt.Println("2")
		}
		fmt.Println("1")
	}
	return nil, err
}
