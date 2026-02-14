package service

import (
	"context"
	"errors"
	"testing"

	"github.com/osvathbotond/cloudant-airportdb-go/internal/model"
)

type mockRepository struct {
	hubs      []model.Hub
	returnErr error
}

func (m *mockRepository) GetByBounds(_ context.Context, minLat, maxLat, minLon, maxLon float64) ([]model.Hub, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	var filtered []model.Hub
	for _, h := range m.hubs {
		if h.Lat < minLat || h.Lat > maxLat {
			continue
		}
		if minLon <= maxLon {
			if h.Lon < minLon || h.Lon > maxLon {
				continue
			}
		} else {
			if h.Lon < minLon && h.Lon > maxLon {
				continue
			}
		}
		filtered = append(filtered, h)
	}
	return filtered, nil
}

func TestNewFinder(t *testing.T) {
	repo := &mockRepository{}
	finder := NewFinder(repo)

	if finder == nil {
		t.Fatal("NewFinder returned nil")
	}

	if finder.repo != repo {
		t.Error("Finder repository not set correctly")
	}
}

func TestFindNearby_EmptyRepository(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestFindNearby_SingleHubWithinRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{
				ID:   "hub1",
				Name: "JFK Airport",
				Lat:  40.6413,
				Lon:  -73.7781,
			},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "hub1" {
		t.Errorf("Expected hub1, got %s", results[0].ID)
	}

	if results[0].DistanceKm <= 0 || results[0].DistanceKm > 50 {
		t.Errorf("Expected distance between 0 and 50 km, got %f", results[0].DistanceKm)
	}
}

func TestFindNearby_SingleHubOutsideRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{
				ID:   "hub1",
				Name: "Los Angeles Airport",
				Lat:  34.0522,
				Lon:  -118.2437,
			},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestFindNearby_MultipleHubsSomeWithinRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{
				ID:   "hub1",
				Name: "JFK Airport",
				Lat:  40.6413,
				Lon:  -73.7781,
			},
			{
				ID:   "hub2",
				Name: "Newark Airport",
				Lat:  40.6895,
				Lon:  -74.1745,
			},
			{
				ID:   "hub3",
				Name: "Los Angeles Airport",
				Lat:  34.0522,
				Lon:  -118.2437,
			},
			{
				ID:   "hub4",
				Name: "LaGuardia Airport",
				Lat:  40.7769,
				Lon:  -73.8740,
			},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	for _, result := range results {
		if result.ID == "hub3" {
			t.Error("Los Angeles Airport should not be in results (too far)")
		}
	}
}

func TestFindNearby_SortingByDistance(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{
				ID:   "hub1",
				Name: "Far Hub",
				Lat:  40.5,
				Lon:  -74.5,
			},
			{
				ID:   "hub2",
				Name: "Close Hub",
				Lat:  40.72,
				Lon:  -74.01,
			},
			{
				ID:   "hub3",
				Name: "Medium Hub",
				Lat:  40.6,
				Lon:  -74.2,
			},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 100)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	for i := 1; i < len(results); i++ {
		if results[i].DistanceKm < results[i-1].DistanceKm {
			t.Errorf("Results not sorted: result[%d] distance (%f) < result[%d] distance (%f)",
				i, results[i].DistanceKm, i-1, results[i-1].DistanceKm)
		}
	}

	if results[0].ID != "hub2" {
		t.Errorf("Expected closest hub to be hub2, got %s", results[0].ID)
	}
}

func TestFindNearby_AllHubsWithinRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Hub 1", Lat: 40.71, Lon: -74.00},
			{ID: "hub2", Name: "Hub 2", Lat: 40.72, Lon: -74.01},
			{ID: "hub3", Name: "Hub 3", Lat: 40.70, Lon: -74.02},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 10)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestFindNearby_AllHubsOutsideRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Hub 1", Lat: 50.0, Lon: -80.0},
			{ID: "hub2", Name: "Hub 2", Lat: 35.0, Lon: -120.0},
			{ID: "hub3", Name: "Hub 3", Lat: 30.0, Lon: -90.0},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 10)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestFindNearby_ZeroRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Exact location", Lat: 40.7128, Lon: -74.0060},
			{ID: "hub2", Name: "Close location", Lat: 40.7129, Lon: -74.0061},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 0)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) > 1 {
		t.Errorf("Expected 0-1 results for zero radius, got %d", len(results))
	}
}

func TestFindNearby_VeryLargeRadius(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "NYC", Lat: 40.7128, Lon: -74.0060},
			{ID: "hub2", Name: "LA", Lat: 34.0522, Lon: -118.2437},
			{ID: "hub3", Name: "London", Lat: 51.5074, Lon: -0.1278},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 10000)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestFindNearby_BoundaryCase(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Boundary Hub", Lat: 40.26, Lon: -74.0},
		},
	}
	finder := NewFinder(repo)

	results1, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50.5)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(results1) != 1 {
		t.Errorf("Expected 1 result with radius 50.5 km, got %d", len(results1))
	}

	results2, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 49.5)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(results2) != 0 {
		t.Errorf("Expected 0 results with radius 49.5 km, got %d", len(results2))
	}
}

func TestFindNearby_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database connection failed")
	repo := &mockRepository{
		returnErr: expectedErr,
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if results != nil {
		t.Errorf("Expected nil results on error, got %d results", len(results))
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error to wrap original error")
	}
}

func TestFindNearby_NegativeCoordinates(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Sydney", Lat: -33.8688, Lon: 151.2093},
			{ID: "hub2", Name: "Near Sydney", Lat: -33.8, Lon: 151.0},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), -33.8688, 151.2093, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected at least 1 result")
	}

	if results[0].ID != "hub1" {
		t.Errorf("Expected closest hub to be hub1, got %s", results[0].ID)
	}
}

func TestFindNearby_DateLineProximity(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "West of date line", Lat: 0.0, Lon: 179.5},
			{ID: "hub2", Name: "East of date line", Lat: 0.0, Lon: -179.5},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 0.0, 180.0, 200)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results near date line, got %d", len(results))
	}
}

func TestFindNearby_PolarRegions(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "North Pole Station", Lat: 89.9, Lon: 0.0},
			{ID: "hub2", Name: "Near North Pole", Lat: 89.5, Lon: 90.0},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 90.0, 0.0, 100)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected at least 1 result near North Pole")
	}
}

func TestFindNearby_HubsSameName(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Airport", Lat: 40.7, Lon: -74.0},
			{ID: "hub2", Name: "Airport", Lat: 40.8, Lon: -74.1},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 50)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0].ID == results[1].ID {
		t.Error("Expected different IDs for different hubs")
	}
}

func TestFindNearby_DistanceCalculation(t *testing.T) {
	repo := &mockRepository{
		hubs: []model.Hub{
			{ID: "hub1", Name: "Test Hub", Lat: 40.7128, Lon: -74.0060},
		},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 10)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].DistanceKm > 0.001 {
		t.Errorf("Expected distance near 0, got %f km", results[0].DistanceKm)
	}
}

func TestFindNearby_PreservesHubData(t *testing.T) {
	expectedHub := model.Hub{
		ID:   "test123",
		Name: "Test Airport",
		Lat:  40.7128,
		Lon:  -74.0060,
	}

	repo := &mockRepository{
		hubs: []model.Hub{expectedHub},
	}
	finder := NewFinder(repo)

	results, err := finder.FindNearby(context.Background(), 40.7128, -74.0060, 10)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	result := results[0]
	if result.ID != expectedHub.ID {
		t.Errorf("Expected ID %s, got %s", expectedHub.ID, result.ID)
	}
	if result.Name != expectedHub.Name {
		t.Errorf("Expected Name %s, got %s", expectedHub.Name, result.Name)
	}
	if result.Lat != expectedHub.Lat {
		t.Errorf("Expected Lat %f, got %f", expectedHub.Lat, result.Lat)
	}
	if result.Lon != expectedHub.Lon {
		t.Errorf("Expected Lon %f, got %f", expectedHub.Lon, result.Lon)
	}
}
