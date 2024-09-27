package main

import (
	"os"
	"time"

	xqueue "github.com/avila-r/message-broker"
	"github.com/gofiber/fiber/v2"
)

var (
	app_url    = os.Getenv("SERVER_URL")
	queue_url  = os.Getenv("QUEUE_URL")
	queue_name = os.Getenv("QUEUE_NAME")
)

func main() {
	// Wait for 5 seconds to assert that RabbitMQ is up
	time.Sleep(5 * time.Second)

	app := fiber.New()

	// Connect to RabbitMQ
	rabbitmq, err := xqueue.NewClient(queue_url)

	if err != nil {
		panic(err)
	}

	// Declare new queue
	rabbitmq.NewQueue(&xqueue.QueueOptions{
		Name: queue_name,
	})

	app.Post("/", func(c *fiber.Ctx) error {
		data := struct{ Message string }{}

		if err := c.BodyParser(&data); err != nil {
			return err
		}

		// Produce a message
		if err := rabbitmq.Produce(xqueue.Message{
			QueueName: queue_name,
			Body:      data,
		}); err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"message": "successful at sending message to queue",
		})
	})

	app.Listen(app_url)
}
