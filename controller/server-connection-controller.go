package controller

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/minhmannh2001/sms/database"
	elastic "github.com/olivere/elastic/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/matryer/try.v1"
)

var serverDatabase database.SMSDatabase = database.NewSMSDatabase()

type ServerHeartBeatResponse struct {
	Ipv4 string    `json:"ipv4"`
	Time time.Time `json:"timestamp"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	}
}`

type ConnectionController interface {
	SendCheckConnection()
	ReceiveCheckConnection()
}

type connectionController struct {
	connection *amqp.Connection
}

func NewConnectionController() ConnectionController {
	var conn *amqp.Connection
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			time.Sleep(10 * time.Second) // wait a minute
		}
		return attempt < 5, err
	})
	// conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("Failed to connect to the rabbitmq")
		return nil
	}
	log.Println("Connected successfully to the rabbitmq")

	return &connectionController{
		connection: conn,
	}
}

func (controller *connectionController) SendCheckConnection() {
	ch, err := controller.connection.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel")
	}
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

	if err != nil {
		log.Fatalln("Failed to declare an exchange")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "sms:ping"
	for {
		err = ch.PublishWithContext(ctx,
			"logs", // exchange
			"",     // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Fatalln("Failed to publish a message")
		}
		log.Printf(" [x] Sent %s", body)
		servers_for_update, err := serverDatabase.ViewServers(0, 0, 0, "", "", "")
		if err != nil {
			panic(err)
		}
		for _, server := range servers_for_update {
			server.Status = "Down"
			serverDatabase.UpdateServer(&server)
		}
		time.Sleep(60 * time.Second)
	}

}

func (controller *connectionController) ReceiveCheckConnection() {
	ch, err := controller.connection.Channel()
	if err != nil {
		log.Fatalln("Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalln("Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalln("Failed to register a consumer")
	}

	ctx := context.Background()

	var client *elastic.Client
	err = try.Do(func(attempt int) (bool, error) {
		var err error
		client, err = elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
		if err != nil {
			time.Sleep(10 * time.Second) // wait a minute
		}
		return attempt < 5, err
	})
	// client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("myindex").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("myindex").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ip := strings.Split(string(d.Body), ":")[0]
			log.Println("ip:", ip)

			server_for_update := serverDatabase.GetServerByIp(ip)
			server_for_update.Status = "Up"
			err := serverDatabase.UpdateServer(&server_for_update)
			if err != nil {
				log.Panicln(err)
			}
			// Create a new document
			doc := ServerHeartBeatResponse{
				Ipv4: ip,
				Time: time.Now(),
			}

			fmt.Println(doc)
			// Index the document
			_, err = client.Index().
				Index("myindex").
				Type("mytype").
				BodyJson(doc).
				Do(ctx)
			if err != nil {
				panic(err)
			}
		}
	}()

	log.Printf("Waiting for messages...")
	<-forever
}

// https://olivere.github.io/elastic/
// https://www.rabbitmq.com/getstarted.html
