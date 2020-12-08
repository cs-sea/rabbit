package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"

	amqppool "github.com/cs-sea/amqp-pool"
)

func main() {
	pool := amqppool.NewConnectionPool(&amqppool.ConnPoolConfig{
		ConnNum: 3,
		Url:     "amqp://guest:guest@localhost:5672/",
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		conn, err := pool.Get()

		if err != nil {
			log.Fatalln(err)
		}

		channelPool := conn.GetChannelPool()

		ch, err := channelPool.Get()

		if err != nil {
			failOnError(err, "get channel")
		}

		_, err = ch.QueueDeclare(&amqppool.QueueDeclareConfig{
			Name:       "task_queue",
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		})
		if err != nil {
			fmt.Println(err)
		}

		err = ch.Publish(&amqppool.PublishData{
			Exchange:  "",
			Key:       "take_queue",
			Mandatory: false,
			Immediate: false,
			Msg: amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte("hello"),
			},
		})

		channelPool.Release(ch)
		pool.Release(conn)
		//pool.ReleaseChannel(ch)
		failOnError(err, "publish")
	})
	fmt.Println(http.ListenAndServe("localhost:9999", mux))

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
