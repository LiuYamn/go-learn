package main

import (
	"go-learn/redis-demo/redigo"
	"log"
	"time"
)

//Redis 全局redis对象
var Redis redigo.Redigo

func main() {
	Redis = redigo.NewRedisPool("127.0.0.1:6379", "", 100, 100, 200)
	if Redis.TestConn() != nil {
		log.Fatal("redis_err", "Redis connect failed!")
	}
	time.Sleep(time.Second * 3)
}
