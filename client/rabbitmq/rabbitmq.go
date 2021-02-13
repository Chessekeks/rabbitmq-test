package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

const (
	mqConnectionAddress = "amqp://localhost:5672/"
)

type Queue struct {
	Name string
	ch   *amqp.Channel
}

func NewQueue(name string) (*Queue, error) {
	var conn *amqp.Connection
	var err error
	for i := range []int{1, 3, 5} {
		if conn, err = connect(mqConnectionAddress); err == nil {
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	if err != nil {
		err = fmt.Errorf("connect to mq error: %w", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = fmt.Errorf("create mq channel error: %w", err)
		return nil, err
	}

	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		err = fmt.Errorf("create queue error: %w", err)
		return nil, err
	}

	return &Queue{Name: name, ch: ch}, nil
}

func (q *Queue) Publish(contentType string, body string) (err error) {
	err = q.ch.Publish("", q.Name, false, false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(body),
		})
	if err != nil {
		err = fmt.Errorf("publish message error: %w", err)
		return
	}

	return
}

func (q *Queue) Receive() (msg <-chan amqp.Delivery, err error) {
	msgCh, err := q.ch.Consume(q.Name,
		"",
		false, false, false, false, nil)
	if err != nil {
		err = fmt.Errorf("subscribe to consume error: %w", err)
		return
	}

	return msgCh, err
}

func connect(addr string) (*amqp.Connection, error) {
	return amqp.Dial(addr)
}
