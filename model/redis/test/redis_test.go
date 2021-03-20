package test

import (
	"fmt"
	"github.com/whereabouts/chassis/model/redis"
	"github.com/whereabouts/utils/code"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	client, err := redis.Init(redis.Config{
		Addr:      ":6379",
		Password:  "root",
		MaxIdle:   10,
		MaxActive: 50,
	})
	defer client.Close()
	if err != nil {
		log.Fatal(err)
	}
	id := "korbin"
	result := GetEmailCache().AddEmailCode(id, code.Default())
	fmt.Println(result)
	result = GetEmailCache().GetEmailCode(id)
	fmt.Println(result.String())
}
