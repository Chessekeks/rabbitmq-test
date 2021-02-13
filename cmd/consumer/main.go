package main

import (
	"log"
	"os"
	"os/signal"
	"rabbitmq-test/client/rabbitmq"
	"syscall"
)

const (
	defaultQueueName = "super-queue"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	queue, err := rabbitmq.NewQueue(defaultQueueName)
	if err != nil {
		log.Fatal(err)
		return
	}

	msgCh, err := queue.Receive()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		select {
		case msg := <-msgCh:
			log.Println("received message: ", string(msg.Body))
		case <-quit:
			log.Println("catch signal from OS - quiting")
			return
		}
	}
}
