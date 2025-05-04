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

var anomalyThresholds map[string]float64 = map[string]float64{
	"PM2.5": 150,
	"PM10":  180,
	"NO2":   150,
	"SO2":   100,
	"O3":    200,
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

	/*
	   These threshold values correspond to extreme situations and
	   might not correspond to real-life thresholds.
	   These exists for testing purposes.
	   These thresholds are also used in the `auto-test.sh` script

	   "PM2.5" anomaly_min=150;
	   "PM10"  anomaly_min=180;
	   "NO2"   anomaly_min=150;
	   "SO2"   anomaly_min=100;
	   "O3"    anomaly_min=200;
	*/

	if math.Abs(zscore) > 2 || entry.Value >= anomalyThresholds[entry.Pollutant] {
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
			Message:   "Anomaly detected!",
			Latitude:  entry.Latitude,
			Longitude: entry.Longitude,
			Value:     entry.Value,
			Pollutant: entry.Pollutant,
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
