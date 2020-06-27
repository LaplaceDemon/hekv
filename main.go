package main

import (
	"hekv/redis"
)

func main() {
	server := redis.CreateServer()
	server.Run()
}
