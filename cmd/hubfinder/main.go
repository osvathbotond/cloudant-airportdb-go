package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/osvathbotond/cloudant-airportdb-go/internal/finder"
	"github.com/osvathbotond/cloudant-airportdb-go/internal/repository"
)

const (
	baseURL = "https://mikerhodes.cloudant.com"
	db      = "airportdb"
	ddoc    = "view1"
	index   = "geo"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	repo, err := repository.NewCloudantRepository(repository.CloudantConfig{
		BaseURL: baseURL,
		DB:      db,
		Ddoc:    ddoc,
		Index:   index,
	})
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}

	f := finder.New(repo)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("This program finds transport hubs within a specified radius from a given point.")

	lat := readFloatUntilValid(scanner, "latitude", -90.0, 90.0)
	lon := readFloatUntilValid(scanner, "longitude", -180.0, 180.0)
	radiusKm := readFloatUntilValid(scanner, "radius in kilometers", 0, 40075)

	hubs, err := f.FindNearby(ctx, lat, lon, radiusKm)
	if err != nil {
		return fmt.Errorf("find nearby hubs: %w", err)
	}

	fmt.Printf("\nFound %d transport hub(s):\n\n", len(hubs))
	for _, hub := range hubs {
		fmt.Printf("Hub: %s\n", hub.Name)
		fmt.Printf("  Distance: %.2f km\n", hub.DistanceKm)
		fmt.Printf("  Latitude: %.6f\n", hub.Lat)
		fmt.Printf("  Longitude: %.6f\n\n", hub.Lon)
	}

	return nil
}
