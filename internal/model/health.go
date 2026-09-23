package model

type HealthResponse struct {
	Status      string `json:"status"`
	Application string `json:"application"`
	Version     string `json:"version"`
}
