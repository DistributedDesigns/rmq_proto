package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", msg, err)
		log.Fatalf(errorMsg)
		panic(errorMsg)
	}
}

const (
	rabbitConnection = "amqp://guest:guest@localhost:44430/"
)

func main() {

	conn, err := amqp.Dial(rabbitConnection)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"quote_broadcast",  // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Receive all quotes
	err = ch.QueueBind(
		q.Name,            // name
		"#",               // routing key
		"quote_broadcast", // exchange
		false,             // no-wait
		nil,               // args
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		log.Println(" [-] Waiting for broadcasts")

		for d := range msgs {
			// Do more intersting things with the quote here
			log.Printf(" [↓] Update: %s", d.Body)
		}
	}()

	<-forever
}
