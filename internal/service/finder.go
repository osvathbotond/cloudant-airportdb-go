package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/osvathbotond/cloudant-airportdb-go/internal/geo"
	"github.com/osvathbotond/cloudant-airportdb-go/internal/model"
	"github.com/osvathbotond/cloudant-airportdb-go/internal/repository"
)

type Finder struct {
	repo repository.Repository
}

func NewFinder(repo repository.Repository) *Finder {
	return &Finder{repo: repo}
}

// FindNearby finds transport hubs within a specified radius (in kilometers) from a given point.
// It returns a slice of hubs with distances sorted by distance from the given point (closest first).
func (f *Finder) FindNearby(ctx context.Context, lat, lon, radiusKm float64) ([]model.HubWithDistance, error) {
	minLat, maxLat, minLon, maxLon, err := geo.CalculateBoundingBox(lat, lon, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("calculate bounding box: %w", err)
	}

	hubs, err := f.repo.GetByBounds(ctx, minLat, maxLat, minLon, maxLon)
	if err != nil {
		return nil, fmt.Errorf("get hubs by bounds: %w", err)
	}

	nearbyHubs := make([]model.HubWithDistance, 0, len(hubs))
	for _, hub := range hubs {
		distanceKm, err := geo.HaversineDistance(lat, lon, hub.Lat, hub.Lon)
		if err != nil {
			return nil, fmt.Errorf("calculate distance for hub %s: %w", hub.ID, err)
		}
		if distanceKm <= radiusKm {
			nearbyHubs = append(nearbyHubs, model.HubWithDistance{
				Hub:        hub,
				DistanceKm: distanceKm,
			})
		}
	}

	sort.Slice(nearbyHubs, func(i, j int) bool {
		return nearbyHubs[i].DistanceKm < nearbyHubs[j].DistanceKm
	})

	return nearbyHubs, nil
}
