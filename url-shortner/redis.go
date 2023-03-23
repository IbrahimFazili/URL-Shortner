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
const env_secondary_redis = "SECONDARY_REDIS_ADDR"

var g_primary_redis *redis.Client
var g_secondary_redis *redis.Client

const g_cache_duration = 120

func redis_init() {

	g_primary_redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv(env_primary_redis)),
	})

	g_secondary_redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv(env_secondary_redis)),
	})

	time.Sleep(3 * time.Second)
	err := g_primary_redis.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Issue connecting to primary redis - %s", err.Error())
	} else {
		log.Printf("Primary redis ready to accept write requests")
	}

	err = g_secondary_redis.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Issue connecting to secondary redis - %s", err.Error())
	} else {
		log.Printf("Replica redis ready to accept read requests")
	}
}

func redis_put(short string, long string) bool {
	_, err := g_primary_redis.Set(context.Background(), short, long, g_cache_duration*time.Second).Result()
	if err != nil {
		log.Printf("[REDIS]: unable to store (%s:%s) in cache\n", short, long)
		return false
	}
	log.Printf("[REDIS]: successfully put (%s:%s) in cache\n", short, long)
	return true
}

func redis_get(short string) string {
	res, err := g_secondary_redis.Get(context.Background(), short).Result()
	if err != nil {
		log.Printf("[REDIS]: Unable to fetch from cache for (%s)\n", short)
		return ""
	}
	log.Printf("[REDIS]: Fetched long from cache for (%s:%s)\n", short, res)
	return res
}
