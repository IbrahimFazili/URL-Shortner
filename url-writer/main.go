package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis/v9"
)

type Url_Data struct {
	Short string
	Long  string
}

const env_pubsub_channel = "PUBSUB_CHANNEL"

func logger_init() {
	os.Remove("logs/output.log")
	file, err := os.OpenFile("logs/output.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Creating log file for url writer")
}

func main() {

	logger_init()
	connect_db_cluster()
	connect_redis()

	pubsub := g_primary_redis.Subscribe(context.Background(), os.Getenv(env_pubsub_channel))
	channel := pubsub.Channel()

	for msg := range channel {

		// to run concurrently
		go func(sub_msg *redis.Message) {
			var data_endata_entry = &Url_Data{}
			// unmarshal to get short long
			json.Unmarshal([]byte(sub_msg.Payload), data_endata_entry)
			log.Printf("Unmarshalled values: %+v\n", data_endata_entry)

			// add to cassandra
			db_put_long(*data_endata_entry)

			// add to redis
			redis_put(*data_endata_entry)
		}(msg)
	}
}
