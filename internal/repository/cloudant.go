package repository

import (
	"context"
	"fmt"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/osvathbotond/cloudant-airportdb-go/internal/model"
)

// Compile-time check that CloudantRepository implements Repository.
var _ Repository = (*CloudantRepository)(nil)

type CloudantRepository struct {
	service *cloudantv1.CloudantV1
	db      string
	ddoc    string
	index   string
}

type CloudantConfig struct {
	BaseURL string
	DB      string
	Ddoc    string
	Index   string
}

func NewCloudantRepository(cfg CloudantConfig) (*CloudantRepository, error) {
	authenticator, err := core.NewNoAuthAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("create no-auth authenticator: %w", err)
	}

	service, err := cloudantv1.NewCloudantV1(&cloudantv1.CloudantV1Options{
		URL:           cfg.BaseURL,
		Authenticator: authenticator,
	})
	if err != nil {
		return nil, fmt.Errorf("create cloudant client: %w", err)
	}

	return &CloudantRepository{
		service: service,
		db:      cfg.DB,
		ddoc:    cfg.Ddoc,
		index:   cfg.Index,
	}, nil
}

// buildSearchQuery constructs a Cloudant Lucene query string for searching
// hubs within the given geographic bounds.
func buildSearchQuery(minLat, maxLat, minLon, maxLon float64) string {
	if minLon > maxLon {
		return fmt.Sprintf("lat:[%f TO %f] AND (lon:[%f TO 180] OR lon:[-180 TO %f])", minLat, maxLat, minLon, maxLon)
	}
	return fmt.Sprintf("lat:[%f TO %f] AND lon:[%f TO %f]", minLat, maxLat, minLon, maxLon)
}

func (r *CloudantRepository) GetByBounds(ctx context.Context, minLat, maxLat, minLon, maxLon float64) ([]model.Hub, error) {
	query := buildSearchQuery(minLat, maxLat, minLon, maxLon)

	var allHubs []model.Hub
	var bookmark *string
	skippedRows := 0

	for {
		options := &cloudantv1.PostSearchOptions{
			Db:       core.StringPtr(r.db),
			Ddoc:     core.StringPtr(r.ddoc),
			Index:    core.StringPtr(r.index),
			Query:    core.StringPtr(query),
			Limit:    core.Int64Ptr(200),
			Bookmark: bookmark,
		}

		result, _, err := r.service.PostSearchWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("post search: %w", err)
		}

		if result.Rows != nil {
			for _, row := range result.Rows {
				if row.ID == nil || row.Fields == nil {
					skippedRows++
					continue
				}

				lat, latOk := row.Fields["lat"].(float64)
				lon, lonOk := row.Fields["lon"].(float64)
				name, nameOk := row.Fields["name"].(string)

				if latOk && lonOk && nameOk {
					allHubs = append(allHubs, model.Hub{
						ID:   *row.ID,
						Lat:  lat,
						Lon:  lon,
						Name: name,
					})
				} else {
					skippedRows++
				}
			}
		}

		if result.Bookmark == nil || *result.Bookmark == "" || len(result.Rows) == 0 {
			break
		}

		bookmark = result.Bookmark
	}

	return allHubs, nil
}
