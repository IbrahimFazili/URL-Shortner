package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
)

var g_keyspace string
var g_table string

var g_ctx = context.Background()
var g_session *gocql.Session

type Url_Data struct {
	Short string
	Long  string
}

func db_init() {
	cluster := gocql.NewCluster(fmt.Sprintf("%s:9042", os.Getenv("CASSANDRA_CONNECT_POINT")))
	cluster.Consistency = gocql.One
	g_session, _ = cluster.CreateSession()

	g_keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	g_table = os.Getenv("CASSANDRA_TABLE")

	err := g_session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 2}", g_keyspace)).Exec()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Created keyspace", g_keyspace)
	}

	cluster.Keyspace = g_keyspace
	err = g_session.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (short text PRIMARY KEY, long text)", g_keyspace, g_table)).Exec()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Created table", g_table)
	}
}

func db_fetch_long(short string) string {
	var data_entry Url_Data
	query_string := fmt.Sprintf("SELECT short, long from %s.%s WHERE short = ?", g_keyspace, g_table)
	err := g_session.Query(query_string, short).WithContext(g_ctx).Scan(&data_entry.Short, &data_entry.Long)
	if err != nil {
		log.Printf("[DB]: Unable to fetch (%s). Error: %s\n", data_entry.Short, err.Error())
	}
	log.Printf("[DB]: fetched from databse (%s:%s)\n", data_entry.Short, data_entry.Long)
	return data_entry.Long
}

func db_put_long(short string, long string) bool {
	query_string := fmt.Sprintf("INSERT INTO %s.%s (short, long) VALUES (?, ?)", g_keyspace, g_table)
	err := g_session.Query(query_string, short, long).Exec()
	if err != nil {
		log.Printf("[DB]: unable to set in databse (%s:%s). Error: %s\n", short, long, err.Error())
		return false
	}
	log.Printf("[DB]: set in databse (%s:%s)", short, long)
	return true
}
