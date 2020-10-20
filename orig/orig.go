package main

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
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

	log.Println("Sending messages")
	for i := 1; i <= 3; i++ {
		publishMessage(ch, "my.o", []byte("MSG_"+strconv.Itoa(i)))
		log.Printf("Sent out: %d", i)
		time.Sleep(3 * time.Second)
	}
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
