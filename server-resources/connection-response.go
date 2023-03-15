package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/matryer/try.v1"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	out, err := exec.Command("hostname", "-I").Output()
	failOnError(err, "Failed to execute os command")

	ips := string(out)
	ip := strings.Split(ips, " ")[0]
	log.Println("ip: " + ip)

	var conn *amqp.Connection
	err = try.Do(func(attempt int) (bool, error) {
		var err error
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			time.Sleep(10 * time.Second) // wait a minute
		}
		return attempt < 5, err
	})

	// conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
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

	// Setup channel to send response
	response_q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var forever chan struct{}

	go func() {
		for d := range msgs {
			resp, err := http.Get("http://sms:8080/checkExistence/" + ip)
			failOnError(err, "Cannot send request to check server's existence")
			body, err := io.ReadAll(resp.Body)
			failOnError(err, "Cannot read the body from response")
			result := string(body)

			if strings.Contains(result, "true") {
				err = ch.PublishWithContext(ctx,
					"",              // exchange
					response_q.Name, // routing key
					false,           // mandatory
					false,           // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(fmt.Sprintf("%s:pong", ip)),
					})
				failOnError(err, "Failed to publish a message")
				_ = d.Body
				log.Printf("%s:pong", ip)
			}
		}
	}()

	log.Printf("Waiting for messages...")
	<-forever
}
