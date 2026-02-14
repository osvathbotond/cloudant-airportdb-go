package repository

import (
	"context"

	"github.com/osvathbotond/cloudant-airportdb-go/internal/model"
)

// Repository defines the interface for retrieving transport hubs
type Repository interface {
	// GetByBounds retrieves all hubs within the specified geographic bounds
	GetByBounds(ctx context.Context, minLat, maxLat, minLon, maxLon float64) ([]model.Hub, error)
}
