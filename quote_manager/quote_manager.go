package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
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
)

var (
	pendingQuoteReqs = make(chan string)
)

type quote struct {
	userID string
	stock  string
	price  float32
}

func (q quote) String() string {
	return fmt.Sprintf("%s,%s,%.2f", q.userID, q.stock, q.price)
}

func handleQuoteBroadcast() {

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

	forever := make(chan bool)

	go func() {
		log.Println(" [-] Waiting for new pending quotes")

		for req := range pendingQuoteReqs {
			go generateAndPublishQuote(req, ch)
		}
	}()

	<-forever
}

func generateAndPublishQuote(req string, ch *amqp.Channel) {
	log.Println(" [.] New pending quote request")
	resp := generateQuote(req)
	log.Printf(" [.] Got a response: %+v", resp)

	err := ch.Publish(
		"quote_broadcast", // exchange
		resp.stock,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: resp.userID,
			Body:          []byte(resp.String()),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [↑] Broadcast update for %s", resp.stock)
}

func generateQuote(s string) quote {
	// assume this parses nicely
	request := strings.Split(s, ",")

	var delayPeriod time.Duration
	if request[1] == "SLOW" {
		// Always give a slow response for this stock
		delayPeriod = time.Second * 20
	} else {
		// Inject a random 0->3 sec delay
		delayPeriod = time.Second * time.Duration(rand.Intn(4))
	}
	delayTimer := time.NewTimer(delayPeriod)
	log.Printf(" [-] Waiting for %.0f sec", delayPeriod.Seconds())
	<-delayTimer.C

	return quote{
		userID: request[0],
		stock:  request[1],
		price:  1000 * rand.Float32(),
	}
}

func handleQuoteRequest() {

	conn, err := amqp.Dial(rabbitConnection)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"quote_req", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// only pull a single message out of the queue at a time
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		log.Println(" [-] Monitoring quote_req queue")

		for d := range msgs {
			log.Printf(" [↓] Received a quote request: %s", d.Body)
			pendingQuoteReqs <- string(d.Body)
			d.Ack(false)
		}
	}()

	<-forever
}

func main() {
	rand.Seed(time.Now().Unix())

	go handleQuoteBroadcast()
	go handleQuoteRequest()

	forever := make(chan bool)
	<-forever
}
