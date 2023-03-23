package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

var g_session *gocql.Session

var g_keyspace string
var g_table string

func connect_db_cluster() {
	cluster := gocql.NewCluster(fmt.Sprintf("%s:9042", os.Getenv("CASSANDRA_CONNECT_POINT")))
	cluster.ConnectTimeout = time.Second * 10
	g_session, _ = cluster.CreateSession()
	g_keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	g_table = os.Getenv("CASSANDRA_TABLE")
}

func db_put_long(data_entry Url_Data) bool {
	query_string := fmt.Sprintf("INSERT INTO %s.%s (short, long) VALUES (?, ?)", g_keyspace, g_table)
	err := g_session.Query(query_string, data_entry.Short, data_entry.Long).Exec()
	if err != nil {
		log.Printf("[DB]: unable to set in databse (%s:%s). Error: %s", data_entry.Short, data_entry.Long, err.Error())
		return false
	}
	log.Printf("[DB]: set in databse (%s:%s)", data_entry.Short, data_entry.Long)
	return true
}
