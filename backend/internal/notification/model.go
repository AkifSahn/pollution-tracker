package notification

type Notification struct {
	Type      int     `json:"type"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Value     float64 `json:"value"`
	Pollutant string  `json:"pollutant"`
}
