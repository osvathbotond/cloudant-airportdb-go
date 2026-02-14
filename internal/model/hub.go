package model

type Hub struct {
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Name string  `json:"name"`
}

type HubWithDistance struct {
	Hub
	DistanceKm float64 `json:"distance_km"`
}
