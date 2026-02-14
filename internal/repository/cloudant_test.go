package repository

import "testing"

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name                           string
		minLat, maxLat, minLon, maxLon float64
		expected                       string
	}{
		{
			name:     "normal bounds",
			minLat:   40.0,
			maxLat:   41.0,
			minLon:   -75.0,
			maxLon:   -73.0,
			expected: "lat:[40.000000 TO 41.000000] AND lon:[-75.000000 TO -73.000000]",
		},
		{
			name:     "date line wrap - minLon greater than maxLon",
			minLat:   -10.0,
			maxLat:   10.0,
			minLon:   170.0,
			maxLon:   -170.0,
			expected: "lat:[-10.000000 TO 10.000000] AND (lon:[170.000000 TO 180] OR lon:[-180 TO -170.000000])",
		},
		{
			name:     "zero-size box",
			minLat:   0.0,
			maxLat:   0.0,
			minLon:   0.0,
			maxLon:   0.0,
			expected: "lat:[0.000000 TO 0.000000] AND lon:[0.000000 TO 0.000000]",
		},
		{
			name:     "full latitude range",
			minLat:   -90.0,
			maxLat:   90.0,
			minLon:   -180.0,
			maxLon:   180.0,
			expected: "lat:[-90.000000 TO 90.000000] AND lon:[-180.000000 TO 180.000000]",
		},
		{
			name:     "negative longitude range no wrap",
			minLat:   50.0,
			maxLat:   55.0,
			minLon:   -5.0,
			maxLon:   5.0,
			expected: "lat:[50.000000 TO 55.000000] AND lon:[-5.000000 TO 5.000000]",
		},
		{
			name:     "date line wrap near boundary",
			minLat:   -1.0,
			maxLat:   1.0,
			minLon:   179.5,
			maxLon:   -179.5,
			expected: "lat:[-1.000000 TO 1.000000] AND (lon:[179.500000 TO 180] OR lon:[-180 TO -179.500000])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildSearchQuery(tt.minLat, tt.maxLat, tt.minLon, tt.maxLon)
			if result != tt.expected {
				t.Errorf("got:  %s\nwant: %s", result, tt.expected)
			}
		})
	}
}
