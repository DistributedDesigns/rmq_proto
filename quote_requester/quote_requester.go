package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

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
	userID           = "jappleseed"
	quoteReqQ        = "quote_req"
	quoteBroadcastQ  = "quote_broadcast"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

func initRMQ() {
	var err error
	conn, err = amqp.Dial(rabbitConnection)
	failOnError(err, "Failed to connect to RabbitMQ")
	// closed in main()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	// closed in main()

	// Ensure all of our standard queues exist.
	// For sending requests
	_, err = ch.QueueDeclare(
		quoteReqQ, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// For listening to quote updates
	err = ch.ExchangeDeclare(
		quoteBroadcastQ,    // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to declare an exchange")
}

func stocksFrom(args []string) []string {
	var s []string
	if (len(args) < 2) || os.Args[1] == "" {
		s = append(s, "AAPL")
	} else {
		s = args[1:]
	}
	return s
}

func requestQuote(stock string, ready <-chan bool) {

	// hold for update watcher
	<-ready
	log.Println(" [↑] Requesting quote for", stock)
	err := ch.Publish(
		"",        // exchange
		quoteReqQ, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(userID + "," + stock),
		})
	failOnError(err, "Failed to publish a message")
}

func watchForQuoteUpdate(stock string, updates chan<- string, ready chan<- bool) {

	// Anonymous Q that filters for updates to stock
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,          // name
		stock,           // routing key
		quoteBroadcastQ, // exchange
		false,           // no-wait
		nil,             // args
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

	log.Println(" [-] Waiting for updates to", stock)
	ready <- true
	for d := range msgs {
		if d.CorrelationId == userID {
			log.Printf(" [↓] Received: %s", d.Body)
		} else {
			log.Printf(" [↙] Intercepted: %s", d.Body)
		}
		break
	}

	updates <- stock
}

func main() {
	rand.Seed(time.Now().Unix())

	initRMQ()

	stocks := stocksFrom(os.Args)
	log.Println(" [.] Getting updates for", stocks)

	updates := make(chan string, len(stocks))
	for _, stock := range stocks {
		ready := make(chan bool, 1)
		go watchForQuoteUpdate(stock, updates, ready)
		go requestQuote(stock, ready)
	}

	for a := 0; a < len(stocks); a++ {
		log.Println(" [x] Got update for", <-updates)
	}

	ch.Close()
	conn.Close()
}
