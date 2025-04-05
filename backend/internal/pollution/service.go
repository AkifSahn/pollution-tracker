package pollution

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/AkifSahn/pollution-tracker/internal/notification"
	"github.com/AkifSahn/pollution-tracker/internal/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PollutionService struct {
	repo PollutionRepo
}

func NewPollutionService(repo PollutionRepo) *PollutionService {
	return &PollutionService{repo: repo}
}

func (s *PollutionService) ProcessAndInsertPollutionEntry(ctx context.Context, entry Pollution) error {
	fromTime := entry.Time.Add(-24 * time.Hour)
	// Get mean and stddev for pollution values for 25 km radius
	mean, stddev, err := s.repo.GetMeanAndStd(ctx, entry.Pollutant, 25, entry.Latitude, entry.Longitude, fromTime, entry.Time)
	if err != nil {
		return fmt.Errorf("failed to get mean and std: %s", err.Error())
	}

	var zscore float64
	if stddev > 0 {
		zscore = (entry.Value - mean) / stddev
	} else {
		zscore = 0
	}

	if math.Abs(zscore) > 2 {
		entry.IsAnomaly = true
	} else {
		entry.IsAnomaly = false
	}

	fmt.Printf("Val: %f, Mean: %f, Std: %f, Z: %f", entry.Value, mean, stddev, zscore)

	if err := s.repo.InsertPollution(ctx, entry); err != nil {
		return fmt.Errorf("failed to insert pollution entry - %s", err.Error())
	}

	if entry.IsAnomaly {
		notification := notification.Notification{
			Type:      1,
			Message:   fmt.Sprintf("Anomaly detected for %s: value %.2f", entry.Pollutant, entry.Value),
			Latitude:  entry.Latitude,
			Longitude: entry.Longitude,
		}

		msg, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Failed to marshal anomaly notification - %s", err.Error())
		} else {
			err = rabbitmq.AmqpCh.Publish(
				"",
				"notification_queue",
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				},
			)
			if err != nil {
				log.Printf("Failed to publish anomaly notification - %s", err.Error())
			}
		}
	}

	return nil
}
