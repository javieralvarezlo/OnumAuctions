package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendSearchAuctions(params []byte) []byte {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer")

	correlationId := fmt.Sprintf("%d", time.Now().Nanosecond())

	err = ch.Publish(
		"",
		"search_auctions",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(params),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if correlationId == d.CorrelationId {
			return d.Body
		}
	}

	return nil
}
