package notification

import (
	"encoding/json"
	"log"

	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
)

func ListenAndConsumeNotifications(hub *Hub) {

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
		// TODO: validate the incoming data
		var notification Notification
		err := json.Unmarshal(d.Body, &notification)
		if err != nil {
			log.Printf("Failed to unmarshal the data - %s", err.Error())
			continue
		}

		hub.broadcast <- d.Body
	}
}
