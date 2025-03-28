package pollution

import "time"

type Pollution struct {
	Time      time.Time `json:"time"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Region    string    `json:"region"`
	Value     float64   `json:"value"`
	IsAnomaly bool      `json:"is_anomaly"`
	Pollutant string    `json:"pollutant"`
}
