package main

import (
	"log"
	"math/rand"
	"rabbitmq-test/client/rabbitmq"
	"strconv"
	"time"
)

const (
	defaultQueueName = "super-queue"
	counterLimit     = 25
	contentType      = "text/plain"
)

func main() {
	queue, err := rabbitmq.NewQueue(defaultQueueName)
	if err != nil {
		log.Fatal(err)
		return
	}

	rr := rand.New(rand.NewSource(time.Now().UnixNano()))

	counter := rr.Intn(counterLimit)
	for i := 0; i < counter; i++ {
		err = queue.Publish(contentType, strconv.Itoa(i))
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("published message ", i)
	}
}
