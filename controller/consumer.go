package controller

import (
	"encoding/json"
	"log"
	"os"

	"github.com/FakJeongTeeNhoi/report-service/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func StartConsumeDataFromQueue(Exchange string, key []string) {

	conn, err := amqp.Dial(os.Getenv("AMQP_URI"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		Exchange, // name
		"topic",  // type
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

	for _, s := range key {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, Exchange, s)
		err = ch.QueueBind(
			q.Name,   // queue name
			s,        // routing key
			Exchange, // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			var reserve model.Reserve
			err := json.Unmarshal(msg.Body, &reserve)
			if err != nil {
				log.Println("Error unmarshalling JSON:", err)
				continue
			}
			log.Printf("Received a message: %s", reserve)
			model.AddReportFromReserve(reserve)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
