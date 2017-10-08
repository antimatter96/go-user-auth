// Package cache deals with caching of similar instances of a type of instance
//
// Redis is used to
package cache

import (
	"fmt"
	"time"

	"../constants"
	"github.com/garyburd/redigo/redis"
)

// pool is the main pool used by all other functions
var pool *redis.Pool

func init() {
	pool = newPool(constants.RedisAddress)
}

// newPool generates a common pool from which we can access connections
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

// GetRoleARNFromCache is used by awsOnboarding.awsRegister to get roleArn from cache
func GetRoleARNFromCache(UUID string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	_, errPingRedis := conn.Do("PING") //REDIS IS DOWN
	if errPingRedis != nil {
		fmt.Println(errPingRedis)
		return "-", errPingRedis
	}

	resRedis, errGetRedis := redis.String(conn.Do("GET", UUID))
	if errGetRedis != nil {
		if errGetRedis.Error() != "redigo: nil returned"{
			fmt.Println("Redis Get Error : ", errGetRedis)
		}
		return "-", errGetRedis
	}

	return resRedis, nil
}
