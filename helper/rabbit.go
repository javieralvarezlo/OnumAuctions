package helper

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func newChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	return ch
}

func SendRPC(name string, data []byte) []byte {
	ch := newChannel()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	FailOnError(err, "Failed to register a consumer")

	correlationId := fmt.Sprintf("%d", time.Now().Nanosecond())
	err = ch.Publish(
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       q.Name,
			Body:          []byte(data),
		})
	FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		if correlationId == d.CorrelationId {
			return d.Body
		}
	}
	return nil
}

func SendDelayedMsg(data []byte, delay int64) {
	ch := newChannel()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		"delayed_exchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{"x-delayed-type": "direct"},
	)
	FailOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"delayed_exchange",
		"delayed_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
			Headers: amqp.Table{
				"x-delay": int32(delay),
			},
		})

	FailOnError(err, "Error publishing delayed message")
}

func RecieveMessages(conn *amqp.Connection, queueName string, processFunction func(amqp.Delivery, *amqp.Channel)) {
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Recieved a message from %s: %s", queueName, d.Body)
		processFunction(d, ch)
	}
}

func RecieveDelayMessages(conn *amqp.Connection, queueName string, processFunction func(amqp.Delivery, *amqp.Channel)) {
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
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
	FailOnError(err, "Error creating exchange")

	q, err := ch.QueueDeclare(
		"delayed_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,             // nombre de la cola
		"delayed_key",      // routing key
		"delayed_exchange", // exchange
		false,
		nil,
	)
	FailOnError(err, "Error binding the exchange to the queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Recieved a delayed message from %s: %s", queueName, d.Body)
		processFunction(d, ch)
	}
}

func ReplyMessage(channel *amqp.Channel, d amqp.Delivery, data []byte) {
	err := channel.Publish(
		"",
		d.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          []byte(data),
			CorrelationId: d.CorrelationId,
		},
	)
	FailOnError(err, "Error sending return message")
}
