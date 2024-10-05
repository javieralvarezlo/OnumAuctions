package main

import (
	"encoding/json"
	"fmt"
	"time"

	h "auctions/helper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendCreationAuction(auction h.Auction) h.Auction {
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
	h.FailOnError(err, "Failed to publish a message")

	var newAuction h.Auction

	for d := range msgs {
		if correlationId == d.CorrelationId {
			json.Unmarshal(d.Body, &newAuction)
			defer sendClosingAuction(newAuction)
			return newAuction
		}
	}
	return newAuction
}

func sendClosingAuction(auction h.Auction) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	h.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	h.FailOnError(err, "Failed to open a channel")
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
	h.FailOnError(err, "Failed to declare an exchange")

	auctionJson, _ := json.Marshal(auction)
	currentTimestamp := time.Now().UnixMilli()
	delay := (auction.BidEndTime - currentTimestamp)
	fmt.Println(111)
	fmt.Printf("EndTime: %d, Ts: %d, Delay: %d", auction.BidEndTime, currentTimestamp, delay)
	fmt.Println(delay)

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

	h.FailOnError(err, "Error publishing delayed message")

}
