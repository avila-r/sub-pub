package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	xqueue "github.com/avila-r/message-broker"
	"github.com/gofiber/fiber/v2"
)

var (
	app_url    = os.Getenv("SERVER_URL")
	queue_url  = os.Getenv("QUEUE_URL")
	queue_name = os.Getenv("QUEUE_NAME")
	database   = struct {
		Data []any `json:"data"`
	}{}
)

func main() {
	// Wait for 10 seconds to assert that RabbitMQ is up
	time.Sleep(10 * time.Second)

	app := fiber.New()

	// Connect to RabbitMQ
	rabbitmq, err := xqueue.NewClient(queue_url)

	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(database)
	})

	// Create new consumer and listen
	// by messages in created queue
	channel, _ := rabbitmq.NewConsumer(&xqueue.Listener{
		QueueName: queue_name,
	})

	// Create goroutine to wait for channel receiveing messages
	go func() {
		// Receive messages from consumer channel
		for message := range channel {
			var (
				payload any
			)

			if err := json.Unmarshal(message.Body, &payload); err != nil {
				log.Printf("error while unmarshalling - %v", err)
			}

			database.Data = append(database.Data, payload)
		}
	}()

	app.Listen(app_url)
}
