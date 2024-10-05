package main

import (
	"encoding/json"
	"fmt"
	"time"

	h "auctions/helper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendCreationBid(bid h.Bid) h.Bid {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	h.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	h.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	h.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	h.FailOnError(err, "Failed to register a consumer")

	correlationId := fmt.Sprintf("%d", time.Now().Nanosecond())
	bidJson, _ := json.Marshal(bid)

	err = ch.Publish(
		"",
		"create_bids",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(bidJson),
		})
	h.FailOnError(err, "Failed to publish a message")
	var newBid h.Bid

	for d := range msgs {
		if correlationId == d.CorrelationId {
			json.Unmarshal(d.Body, &newBid)
			return newBid
		}
	}

	return newBid
}
