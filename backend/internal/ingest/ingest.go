package ingest

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/AkifSahn/pollution-tracker/internal/database"
	"github.com/AkifSahn/pollution-tracker/internal/pollution"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
)

func ListenIngestion() {
	msgs, err := rabbitmq.AmqpCh.Consume(
		"ingest_queue", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err.Error())
	}

	repo := pollution.NewPollutionRepo(database.DB)
	service := pollution.NewPollutionService(repo)
	go func() {
		for d := range msgs {
			var data pollution.Pollution
			// TODO: validate the data before unmarshaling, do the validation either here or before
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Failed to unmarshal the data - %s", err.Error())
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			// Handle anomaly detection before inserting into the database
			if err = service.ProcessAndInsertPollutionEntry(ctx, data); err != nil {
				log.Printf("Failed to insert the data into database - %s", err.Error())
			}
			cancel()
		}
	}()
}
