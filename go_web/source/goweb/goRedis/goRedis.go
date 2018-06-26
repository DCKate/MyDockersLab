package goRedis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var (
	mutex  sync.Mutex
	client *redis.Client
)

func GetRedis(hostaddr string) *redis.Client {
	if client == nil {
		mutex.Lock()
		if client == nil {
			client = redis.NewClient(&redis.Options{
				Addr:     hostaddr,
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			pong, err := client.Ping().Result()
			fmt.Println(pong, err)
		}
		mutex.Unlock()
	}
	return client
}
