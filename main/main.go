package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	amqppool "github.com/cs-sea/apqp-pool"
)

func main() {
	pool := amqppool.NewPool(5, "amqp://guest:guest@localhost:5672/")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		channel, err := pool.GetChannel()
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		b, err := json.Marshal(channel)
		if err != nil {
			writer.WriteHeader(501)
			return
		}

		fmt.Println(b)
		_, err = writer.Write(b)
		if err != nil {
			writer.WriteHeader(502)
			return
		}
		fmt.Println(b)
	})

	fmt.Println(http.ListenAndServe("localhost:9999", mux))
}
