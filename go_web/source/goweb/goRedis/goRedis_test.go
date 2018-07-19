package goRedis

import (
	"log"
	"testing"
	"time"
)

func TestSaveRedisSet(t *testing.T) {
	err := SaveRedisSet("QWERTYUIOP0987654321", RedisData{
		Target:  "kkk",
		Message: "welcome",
	}, time.Minute)
	if err != nil {
		log.Println("1", err)
	}
	out, err := GetRedisSet("QWERTYUIOP0987654321")
	if err != nil {
		log.Println("2", err)
	}
	log.Println("sss ", out)
}
func TestSaveRedisMap(t *testing.T) {
	err := SaveRedisMap("ASDFGHJKL1234567890", map[string]interface{}{
		"name":  "kkk",
		"level": 30,
	}, time.Minute)
	if err != nil {
		log.Println("1", err)
	}
	out, err := GetRedisMap("ASDFGHJKL1234567890", "name", "level")
	if err != nil {
		log.Println("2", err)
	}
	log.Println("sss ", out)
}
