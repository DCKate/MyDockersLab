package goRedis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	mutex  sync.Mutex
	client *redis.Client
)

type RedisData struct {
	Target  string `json:"target,omitempty"`
	Message string `json:"msg,omitempty"`
}

func SaveRedisSet(rkey string, data RedisData, tt time.Duration) error {
	var tmp bytes.Buffer
	enc := gob.NewEncoder(&tmp)
	err := enc.Encode(data)
	if err != nil {
		return err
	}
	err = GetRedis().SAdd(rkey, tmp.Bytes()).Err()
	if err != nil {
		return err
	}
	return GetRedis().Expire(rkey, tt).Err()
}

func GetRedisSet(rkey string) (RedisData, error) {
	r := RedisData{}
	val, err := GetRedis().SPop(rkey).Bytes()
	if err != nil {
		return r, err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(val))
	err = dec.Decode(&r)
	return r, err
}

func SaveRedisMap(rkey string, data map[string]interface{}, tt time.Duration) error {
	err := GetRedis().HMSet(rkey, data).Err()
	if err != nil {
		return err
	}
	return GetRedis().Expire(rkey, tt).Err()
}

func GetRedisMap(rkey string, keys ...string) (map[string]interface{}, error) {
	val, err := GetRedis().HMGet(rkey, keys...).Result()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{}, len(val))
	for ii, vv := range val {
		data[keys[ii]] = vv
	}
	return data, nil
}

func SaveRedisData(rkey, tar, msg string, tt time.Duration) error {
	r := RedisData{
		Target:  tar,
		Message: msg,
	}
	vv, _ := json.Marshal(r)
	return GetRedis().Set(rkey, string(vv), tt).Err()
}

func GetRedisData(rkey string) (RedisData, error) {
	r := RedisData{}
	val, err := GetRedis().Get(rkey).Result()
	if err != nil {
		return r, err
	}
	err = json.Unmarshal([]byte(val), &r)
	return r, err
}

func GetRedis() *redis.Client {
	if client == nil {
		mutex.Lock()
		if client == nil {
			client = redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
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
