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

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"create_auctions",
		false,
		false,
		false,
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
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
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

		err = ch.Publish(
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

}
