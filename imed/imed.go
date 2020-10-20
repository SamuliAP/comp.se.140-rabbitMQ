package main

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func main() {

	time.Sleep(10*time.Second)

	conn, err := getConnection()
	onError(err, "Connection refused")
	defer conn.Close()

	ch, err := conn.Channel()
	onError(err, "Couldn't open a channel")
	defer ch.Close()

	declareExhange(ch)
	q := declareQueue(ch)
	bindQueue(ch, q, "my.i")
	msgs := consumeChannel(err, ch, q)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			msg := append([]byte("Got "), d.Body[:]...)
			publishMessage(ch, "my.o", msg)
			log.Printf("Published: %s", msg)
		}
	}()

	log.Println("Listening for messages...")

	<-forever
}

func publishMessage(ch *amqp.Channel, key string, message []byte) {
	err := ch.Publish(
		"comps400", // exchange
		key,     // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	onError(err, "Failed to publish a message")
}

func consumeChannel(err error, ch *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	onError(err, "Failed to register a consumer")
	return msgs
}

func declareExhange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"comps400", // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	onError(err, "Failed to declare an exchange")
}

func declareQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	onError(err, "Failed to declare a queue")
	return q
}

func bindQueue(ch *amqp.Channel, q amqp.Queue, key string) {
	err := ch.QueueBind(
		q.Name,       // queue name
		key,       // routing key
		"comps400", // exchange
		false,
		nil)
	onError(err, "Failed to bind a queue")
}

func getConnection() (*amqp.Connection, error) {

	// Connection might not be immediately available,
	// so poll the connection for ~30 seconds
	for i := 0; i <= 30; i++ {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			return conn, nil
		}
		time.Sleep(1 * time.Second)
	}
	return nil, errors.New("url: amqp://guest:guest@rabbitmq:5672/")
}

func onError(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}