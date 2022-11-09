package redis

import (
	"context"
	"fmt"
	"sync"

	goredis "github.com/go-redis/redis/v8"
)

var lock = &sync.Mutex{}

type Singleton struct {
	Ctx    context.Context
	Client *goredis.Client
}

var singleInstance *Singleton

func GetRedisInstance() *Singleton {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating RedisClient instance now.")
			singleInstance = &Singleton{}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}

func (r Singleton) Connect() *Singleton {
	singleInstance.Client = goredis.NewClient(&goredis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	singleInstance.Ctx = context.Background()
	return singleInstance
}
