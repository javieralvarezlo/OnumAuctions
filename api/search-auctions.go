package main

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	h "auctions/helper"
)

func sendSearchAuctions(params h.AuctionSearchParams) []h.Auction {
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
	paramsJson, _ := json.Marshal(params)

	err = ch.Publish(
		"",
		"search_auctions",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(paramsJson),
		})
	h.FailOnError(err, "Failed to publish a message")

	var auctions []h.Auction

	for d := range msgs {
		if correlationId == d.CorrelationId {
			json.Unmarshal(d.Body, &auctions)
			return auctions
		}
	}

	return nil
}
