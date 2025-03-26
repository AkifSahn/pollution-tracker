package models

import "time"

type Pollution struct {
	Time      time.Time `json:"time"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Pollution float64   `json:"pollution"`
}
