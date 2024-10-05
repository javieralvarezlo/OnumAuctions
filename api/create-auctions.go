package main

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendCreationAuction(auction Auction) Auction {
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
	auctionJson, _ := json.Marshal(auction)

	err = ch.Publish(
		"",
		"create_auctions",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(auctionJson),
		})
	failOnError(err, "Failed to publish a message")

	var newAuction Auction

	for d := range msgs {
		if correlationId == d.CorrelationId {
			json.Unmarshal(d.Body, &newAuction)
			defer sendClosingAuction(newAuction)
			return newAuction
		}
	}
	return newAuction
}

func sendClosingAuction(auction Auction) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	)
	failOnError(err, "Failed to declare an exchange")

	auctionJson, _ := json.Marshal(auction)
	currentTimestamp := time.Now().Unix()
	delay := (auction.BidEndTime - currentTimestamp) * 1000

	err = ch.Publish(
		"delayed_exchange",
		"delayed_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(auctionJson),
			Headers: amqp.Table{
				"x-delay": int32(delay),
			},
		})

	failOnError(err, "Error publishing delayed message")

}
