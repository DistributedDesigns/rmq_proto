package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	rabbitHost              = "localhost"
	rabbitPort              = "44430"
	rabbitConnectionAddress = "amqp://guest:guest@" + rabbitHost + ":" + rabbitPort
)

func failOnError(err error, msg string) {
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", msg, err)
		log.Fatalf(errorMsg)
		panic(errorMsg)
	}
}

func main() {

	conn, err := amqp.Dial(rabbitConnectionAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
