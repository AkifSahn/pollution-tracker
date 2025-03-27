package rabbitmq

import (
	"fmt"
	"log"

	"github.com/AkifSahn/pollution-tracker/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var AmqpConn *amqp.Connection
var AmqpCh *amqp.Channel

func Connect(cfg *config.Config) *amqp.Channel {
	if AmqpCh == nil {
		connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.AmqpUser, cfg.AmqpPassword, cfg.AmqpHost, cfg.AmqpPort)
		conn, err := amqp.Dial(connStr)
		if err != nil {
			log.Fatalf("Failed to connect RabbitMQ: %s", err.Error())
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err.Error())
		}

		AmqpConn = conn
		AmqpCh = ch
	}

	return AmqpCh
}

func DeclareQueues() {
	// We need two queues
	_, err := AmqpCh.QueueDeclare("ingest_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err.Error())
	}

	_, err = AmqpCh.QueueDeclare("notification_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err.Error())
	}
}
