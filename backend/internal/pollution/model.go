package pollution

import "time"

type Pollution struct {
	Time      time.Time `json:"time"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Value     float64   `json:"value"`
	IsAnomaly bool      `json:"is_anomaly"`
	Pollutant string    `json:"pollutant"`
}

type PollutionDensity struct {
	Time      time.Time `json:"time"`
	Pollutant string    `json:"pollutant"`
	Density   float64   `json:"density"`
}

type PollutionValueResponse struct {
	Time      time.Time `json:"time"`
	Value     float64   `json:"value"`
	Pollutant string    `json:"pollutant"`
}
