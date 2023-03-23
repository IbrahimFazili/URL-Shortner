package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

const g_short = "short"
const g_long = "long"
const env_pubsub_channel = "PUBSUB_CHANNEL"

var g_channel string
var g_writer bool

func handle_server_status(responseWrite http.ResponseWriter, request *http.Request) {
	responseWrite.WriteHeader(http.StatusOK)
	io.WriteString(responseWrite, "Server is alive\n")
}

func handle_redirect_short(responseWrite http.ResponseWriter, request *http.Request) {
	var short = mux.Vars(request)[g_short]
	log.Printf("Got GET request with params short=%s\n", short)

	if len(short) == 0 {
		http.Error(responseWrite, "bad request", http.StatusBadRequest)
	}

	fetched_long := redis_get(short)
	if len(fetched_long) > 0 {
		http.Redirect(responseWrite, request, fetched_long, http.StatusTemporaryRedirect)
		return
	}

	fetched_long = db_fetch_long(short)
	if len(fetched_long) > 0 {
		http.Redirect(responseWrite, request, fetched_long, http.StatusTemporaryRedirect)

		redis_put(short, fetched_long)
		return
	}

	if len(fetched_long) == 0 {
		http.Error(responseWrite, "page not found", http.StatusNotFound)
	}
}

func handle_put_short(responseWrite http.ResponseWriter, request *http.Request) {
	query_parameters := request.URL.Query()
	short_query := query_parameters.Get(g_short)
	long_query := query_parameters.Get(g_long)

	log.Printf("Got PUT request with params short=%s and long=%s\n", short_query, long_query)


	if len(short_query) == 0 || len(long_query) == 0 {
		http.Error(responseWrite, "bad request", http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(long_query)
	if err != nil {
		http.Error(responseWrite, "bad request", http.StatusBadRequest)
		return
	}

	if !g_writer {

		db_put_long(short_query, long_query)

		redis_put(short_query, long_query)
	} else {
		data_endata_entry := Url_Data{
			Short: short_query,
			Long:  long_query,
		}
		serialized, err := json.Marshal(data_endata_entry)
		if err != nil {
			// error handling
		}
		go g_primary_redis.Publish(context.Background(), g_channel, string(serialized))
		log.Printf("Published to \"%s\" channel\n", g_channel)
	}

	responseWrite.WriteHeader(http.StatusOK)
}

func logger_init() {
	os.Remove("logs/output.log")
	file, err := os.OpenFile("logs/output.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Creating log file")
}

func main() {
	g_channel = os.Getenv(env_pubsub_channel)

	router := mux.NewRouter()
	router.HandleFunc("/", handle_put_short).Methods("PUT")
	router.HandleFunc("/{short:[a-zA-Z0-9]+}", handle_redirect_short).Methods("GET")

	logger_init()

	g_writer = os.Getenv("EXTERNAL_WRITER") == "true"

	db_init()
	log.Println("Connected to Cassandra")

	redis_init()
	log.Println("Connected to Redis")

	http.ListenAndServe(":9000", router)
}
