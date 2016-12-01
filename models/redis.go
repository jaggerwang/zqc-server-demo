package models

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

var RedisPools = map[string]*redis.Pool{}

func RedisPool(name string) *redis.Pool {
	if _, ok := RedisPools[name]; !ok {
		prefix := fmt.Sprintf("redis.%s", name)
		RedisPools[name] = &redis.Pool{
			MaxIdle:     viper.GetInt(fmt.Sprintf("%s.%s", prefix, "maxIdle")),
			IdleTimeout: viper.GetDuration(fmt.Sprintf("%s.%s", prefix, "idleTimeout")) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout(
					"tcp",
					viper.GetString(fmt.Sprintf("%s.%s", prefix, "address")),
					viper.GetDuration(fmt.Sprintf("%s.%s", prefix, "timeout.connect"))*time.Second,
					viper.GetDuration(fmt.Sprintf("%s.%s", prefix, "timeout.read"))*time.Second,
					viper.GetDuration(fmt.Sprintf("%s.%s", prefix, "timeout.write"))*time.Second,
				)
				if err != nil {
					return nil, err
				}
				password := viper.GetString(fmt.Sprintf("%s.%s", prefix, "password"))
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}
	return RedisPools[name]
}
