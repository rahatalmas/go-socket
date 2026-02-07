package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GoRabbitProducer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Rabbit Tester",
		false, false, false, false, nil,
	)
	failOnError(err, "Failed to declare queue")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	msgContainer := "hello world from producer"

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msgContainer),
		},
	)
	failOnError(err, "failed to push the message")
	log.Printf(" [x] Sent %s\n", msgContainer)
}
