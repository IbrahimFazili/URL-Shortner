package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
)

const env_primary_redis = "PRIMARY_REDIS_ADDR"

var g_primary_redis *redis.Client

const g_cache_duration = 120

func connect_redis() {

	g_primary_redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv(env_primary_redis)),
	})

	time.Sleep(3 * time.Second)
	// make for loop three times
	err := g_primary_redis.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Issue connecting to primary redis - %s", err.Error())
	} else {
		log.Printf("Primary redis ready to accept write requests")
	}
}

func redis_put(data_entry Url_Data) bool {
	_, err := g_primary_redis.Set(context.Background(), data_entry.Short, data_entry.Long, g_cache_duration*time.Second).Result()
	if err != nil {
		log.Printf("[REDIS]: unable to store (%s:%s) in cache\n", data_entry.Short, data_entry.Long)
		return false
	}
	log.Printf("[REDIS]: successfully put (%s:%s) in cache\n", data_entry.Short, data_entry.Long)
	return true
}
