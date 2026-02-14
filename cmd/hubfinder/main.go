package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"text/tabwriter"

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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tDistance (km)\tLatitude\tLongitude")
	fmt.Fprintln(w, "----\t-------------\t--------\t---------")

	for _, hub := range hubs {
		fmt.Fprintf(w, "%s\t%.2f\t%.6f\t%.6f\n", hub.Name, hub.DistanceKm, hub.Lat, hub.Lon)
	}

	w.Flush()

	return nil
}
