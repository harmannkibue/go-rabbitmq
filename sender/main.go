package main

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang-rabbit-mq/consumer"
	"github.com/streadway/amqp"
	"net/http"
	"os"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	go consumer.Consumer()

	router.GET("/send", func(c *gin.Context) {
		// Create a message to publish.
		body := "{virtual-account:1234567891, client-username:wasoko}"
		message := amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		}

		// Attempt to publish a message to the queue.
		if err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to send message",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "sent successfully",
		})
	})

	router.Run()
}
