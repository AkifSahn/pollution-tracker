package notification

import (
	"encoding/json"
	"log"

	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	"github.com/gofiber/websocket/v2"
)

// TODO: Handling multiple connections???
// Currently we are consuming the notification from rabbitmq and it's gone
// Only one client receives the message

// This handler subscribes to RabbitMQ notifications query and publishes them to the connection
func NotificationHandler(c *websocket.Conn) {
	defer c.Close()

	msgs, err := rabbitmq.AmqpCh.Consume(
		"notification_queue", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err.Error())
	}

	for d := range msgs {
		var data Notification

		// TODO: validate the data before writing to the connection
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			log.Printf("Failed to unmarshal the data - %s", err.Error())
			continue
		}

		err = c.WriteMessage(1, d.Body)
		if err != nil {
			log.Printf("Failed to write WriteMessage to websocket connection - %s", err.Error())
			continue
		}
	}
}
