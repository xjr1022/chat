package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//定义一个全局变量
var pool *redis.Pool

func initPool(address string,maxIdle,maxActive int,idleTimeout time.Duration){
	pool = &redis.Pool{
		MaxIdle:maxIdle,
		MaxActive:maxActive,
		IdleTimeout:idleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp",address)
		},
	}
}