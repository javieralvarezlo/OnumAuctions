package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	h "auctions/helper"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	h.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	go recieveMessages(conn, "create_auctions", handleAuctionCreate)
	go recieveMessages(conn, "search_auctions", handleSearchAuctions)
	go recieveDelayMessages(conn, "close_auctions", handleCloseAuctions)

	select {}
}

func recieveMessages(conn *amqp.Connection, queueName string, processFunction func(amqp.Delivery, *amqp.Channel)) {
	ch, err := conn.Channel()
	h.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	h.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Recieved a message from %s: %s", queueName, d.Body)
		processFunction(d, ch)
	}
}

func recieveDelayMessages(conn *amqp.Connection, queueName string, processFunction func(amqp.Delivery, *amqp.Channel)) {
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
	h.FailOnError(err, "Error creating exchange")

	q, err := ch.QueueDeclare(
		"delayed_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	h.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,             // nombre de la cola
		"delayed_key",      // routing key
		"delayed_exchange", // exchange
		false,
		nil,
	)
	h.FailOnError(err, "Error binding the exchange to the queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Recieved a delayed message from %s: %s", queueName, d.Body)
		processFunction(d, ch)
	}
}

func handleAuctionCreate(d amqp.Delivery, channel *amqp.Channel) {
	var auction h.Auction
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
	h.FailOnError(err, "Error sending return message")
}

func handleSearchAuctions(d amqp.Delivery, channel *amqp.Channel) {
	var params h.AuctionSearchParams
	json.Unmarshal(d.Body, &params)

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
	h.FailOnError(err, "Error sending return message")
}

func handleCloseAuctions(d amqp.Delivery, channel *amqp.Channel) {
	var auction h.Auction
	json.Unmarshal(d.Body, &auction)

	closeAuction(auction)

}
