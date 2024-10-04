package main

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendSearchBids(params BidSearchParams) []Bid {
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
	paramsJson, _ := json.Marshal(params)

	err = ch.Publish(
		"",
		"search_bids",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(paramsJson),
		})
	failOnError(err, "Failed to publish a message")

	var bids []Bid

	for d := range msgs {
		if correlationId == d.CorrelationId {
			json.Unmarshal(d.Body, &bids)
			return bids
		}
	}

	return nil
}
