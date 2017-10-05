package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ!")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel!")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue!")

	scanner := bufio.NewScanner(os.Stdin)
    	var text string

    	for text != "exit" {
		fmt.Print("Enter your text: ")
		scanner.Scan()
		text = scanner.Text()
		
		if text != "exit" {
			err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(text),
			})

			log.Printf(" [x] Sent %s", text)
			failOnError(err, "Failed to publish a message!")
        	}
    	}
}
