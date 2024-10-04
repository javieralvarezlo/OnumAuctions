package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	go recieveMessages(conn, "create_auctions", handleAuctionCreate)
	go recieveMessages(conn, "search_auctions", handleSearchAuctions)

	select {}
}

func recieveMessages(conn *amqp.Connection, queueName string, processFunction func(amqp.Delivery, *amqp.Channel)) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Recieved a message from %s: %s", queueName, d.Body)
		processFunction(d, ch)
	}
}

func handleAuctionCreate(d amqp.Delivery, channel *amqp.Channel) {
	var auction Auction
	err := json.Unmarshal(d.Body, &auction)

	if err != nil {
		fmt.Print("Error Unmarshaling")
	}

	auction = createAuction(auction)
	auctionJson, err := json.Marshal(auction)
	if err != nil {
		fmt.Printf("Error marshaling: %s", err)
	}

	err = channel.Publish(
		"",
		d.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(auctionJson),
			CorrelationId: d.CorrelationId,
		},
	)
	failOnError(err, "Error sending return message")
}

func handleSearchAuctions(d amqp.Delivery, channel *amqp.Channel) {
	var params AuctionSearchParams
	err := json.Unmarshal(d.Body, &params)

	if err != nil {
		fmt.Print("Error Unmarshaling")
	}

	result := getAllAuctions(params)
	resultJson, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshaling: %s", err)
	}

	err = channel.Publish(
		"",
		d.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(resultJson),
			CorrelationId: d.CorrelationId,
		},
	)
	failOnError(err, "Error sending return message")
}
