package geo

import (
	"fmt"
	"math"
	"testing"
)

const floatTolerance = 1e-6

func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestDegToRad(t *testing.T) {
	tests := []struct {
		name     string
		degrees  float64
		expected float64
	}{
		{"Zero degrees", 0, 0},
		{"90 degrees", 90, math.Pi / 2},
		{"180 degrees", 180, math.Pi},
		{"270 degrees", 270, 3 * math.Pi / 2},
		{"360 degrees", 360, 2 * math.Pi},
		{"Negative 90 degrees", -90, -math.Pi / 2},
		{"Negative 180 degrees", -180, -math.Pi},
		{"45 degrees", 45, math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := degToRad(tt.degrees)
			if !floatEquals(result, tt.expected, floatTolerance) {
				t.Errorf("degToRad(%f) = %f, want %f", tt.degrees, result, tt.expected)
			}
		})
	}
}

func TestRadToDeg(t *testing.T) {
	tests := []struct {
		name     string
		radians  float64
		expected float64
	}{
		{"Zero radians", 0, 0},
		{"Pi/2 radians", math.Pi / 2, 90},
		{"Pi radians", math.Pi, 180},
		{"3Pi/2 radians", 3 * math.Pi / 2, 270},
		{"2Pi radians", 2 * math.Pi, 360},
		{"Negative Pi/2", -math.Pi / 2, -90},
		{"Negative Pi", -math.Pi, -180},
		{"Pi/4 radians", math.Pi / 4, 45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := radToDeg(tt.radians)
			if !floatEquals(result, tt.expected, floatTolerance) {
				t.Errorf("radToDeg(%f) = %f, want %f", tt.radians, result, tt.expected)
			}
		})
	}
}

func TestDegToRadAndBack(t *testing.T) {
	testValues := []float64{0, 45, 90, 135, 180, -45, -90, -135, -180, 12.34, -67.89}

	for _, degrees := range testValues {
		t.Run(fmt.Sprintf("%.2f degrees", degrees), func(t *testing.T) {
			radians := degToRad(degrees)
			backToDegrees := radToDeg(radians)
			if !floatEquals(degrees, backToDegrees, floatTolerance) {
				t.Errorf("Round trip failed: %f -> %f -> %f", degrees, radians, backToDegrees)
			}
		})
	}
}

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name             string
		lat1, lon1       float64
		lat2, lon2       float64
		expectedDistance float64
		toleranceKm      float64
	}{
		{
			name:             "Same point",
			lat1:             40.7128,
			lon1:             -74.0060,
			lat2:             40.7128,
			lon2:             -74.0060,
			expectedDistance: 0,
			toleranceKm:      0.001,
		},
		{
			name:             "NYC to LA",
			lat1:             40.7128,
			lon1:             -74.0060,
			lat2:             34.0522,
			lon2:             -118.2437,
			expectedDistance: 3936,
			toleranceKm:      50,
		},
		{
			name:             "London to Paris",
			lat1:             51.5074,
			lon1:             -0.1278,
			lat2:             48.8566,
			lon2:             2.3522,
			expectedDistance: 344,
			toleranceKm:      10,
		},
		{
			name:             "Equator span - 90 degrees longitude",
			lat1:             0.0,
			lon1:             0.0,
			lat2:             0.0,
			lon2:             90.0,
			expectedDistance: 10007.5,
			toleranceKm:      10,
		},
		{
			name:             "North Pole to South Pole",
			lat1:             90.0,
			lon1:             0.0,
			lat2:             -90.0,
			lon2:             0.0,
			expectedDistance: 20015,
			toleranceKm:      20,
		},
		{
			name:             "Crossing International Date Line",
			lat1:             0.0,
			lon1:             179.0,
			lat2:             0.0,
			lon2:             -179.0,
			expectedDistance: 222.6,
			toleranceKm:      5,
		},
		{
			name:             "Very close points",
			lat1:             40.7128,
			lon1:             -74.0060,
			lat2:             40.7138,
			lon2:             -74.0070,
			expectedDistance: 0.13,
			toleranceKm:      0.05,
		},
		{
			name:             "Same latitude different longitude",
			lat1:             45.0,
			lon1:             0.0,
			lat2:             45.0,
			lon2:             10.0,
			expectedDistance: 786,
			toleranceKm:      10,
		},
		{
			name:             "Same longitude different latitude",
			lat1:             40.0,
			lon1:             -74.0,
			lat2:             50.0,
			lon2:             -74.0,
			expectedDistance: 1111,
			toleranceKm:      10,
		},
		{
			name:             "North Pole to equator",
			lat1:             90.0,
			lon1:             0.0,
			lat2:             0.0,
			lon2:             0.0,
			expectedDistance: 10007.5,
			toleranceKm:      10,
		},
		{
			name:             "South Pole to equator",
			lat1:             -90.0,
			lon1:             0.0,
			lat2:             0.0,
			lon2:             0.0,
			expectedDistance: 10007.5,
			toleranceKm:      10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance, err := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !floatEquals(distance, tt.expectedDistance, tt.toleranceKm) {
				t.Errorf("HaversineDistance(%f, %f, %f, %f) = %f km, want %f km (±%f km)",
					tt.lat1, tt.lon1, tt.lat2, tt.lon2, distance, tt.expectedDistance, tt.toleranceKm)
			}
		})
	}
}

func TestHaversineDistanceSymmetry(t *testing.T) {
	testCases := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
	}{
		{"NYC to LA", 40.7128, -74.0060, 34.0522, -118.2437},
		{"London to Paris", 51.5074, -0.1278, 48.8566, 2.3522},
		{"Across date line", 0.0, 179.0, 0.0, -179.0},
		{"North to South Pole", 90.0, 0.0, -90.0, 0.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distAB, err := HaversineDistance(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
			if err != nil {
				t.Fatalf("unexpected error for A->B: %v", err)
			}
			distBA, err := HaversineDistance(tc.lat2, tc.lon2, tc.lat1, tc.lon1)
			if err != nil {
				t.Fatalf("unexpected error for B->A: %v", err)
			}
			if !floatEquals(distAB, distBA, floatTolerance) {
				t.Errorf("Symmetry violated: A->B = %f km, B->A = %f km", distAB, distBA)
			}
		})
	}
}

func TestHaversineDistanceTriangleInequality(t *testing.T) {
	nycLat, nycLon := 40.7128, -74.0060
	lonLat, lonLon := 51.5074, -0.1278
	parLat, parLon := 48.8566, 2.3522

	distNYCLon, err := HaversineDistance(nycLat, nycLon, lonLat, lonLon)
	if err != nil {
		t.Fatalf("unexpected error for NYC->London: %v", err)
	}
	distLonPar, err := HaversineDistance(lonLat, lonLon, parLat, parLon)
	if err != nil {
		t.Fatalf("unexpected error for London->Paris: %v", err)
	}
	distNYCPar, err := HaversineDistance(nycLat, nycLon, parLat, parLon)
	if err != nil {
		t.Fatalf("unexpected error for NYC->Paris: %v", err)
	}

	if distNYCLon+distLonPar < distNYCPar-1.0 {
		t.Errorf("Triangle inequality violated: NYC->Lon(%f) + Lon->Par(%f) < NYC->Par(%f)",
			distNYCLon, distLonPar, distNYCPar)
	}
}

func TestCalculateBoundingBox(t *testing.T) {
	tests := []struct {
		name                           string
		lat, lon, radiusKm             float64
		validateMinLat, validateMaxLat float64
		validateMinLon, validateMaxLon float64
		toleranceDeg                   float64
	}{
		{
			name:           "Normal case - NYC 50km radius",
			lat:            40.7128,
			lon:            -74.0060,
			radiusKm:       50,
			validateMinLat: 40.2632,
			validateMaxLat: 41.1624,
			validateMinLon: -74.599,
			validateMaxLon: -73.413,
			toleranceDeg:   0.02,
		},
		{
			name:           "Zero radius",
			lat:            40.7128,
			lon:            -74.0060,
			radiusKm:       0,
			validateMinLat: 40.7128,
			validateMaxLat: 40.7128,
			validateMinLon: -74.0060,
			validateMaxLon: -74.0060,
			toleranceDeg:   0.001,
		},
		{
			name:           "Point at equator",
			lat:            0.0,
			lon:            0.0,
			radiusKm:       100,
			validateMinLat: -0.8993,
			validateMaxLat: 0.8993,
			validateMinLon: -0.8993,
			validateMaxLon: 0.8993,
			toleranceDeg:   0.02,
		},
		{
			name:           "Near prime meridian",
			lat:            51.5074,
			lon:            0.0,
			radiusKm:       50,
			validateMinLat: 51.0578,
			validateMaxLat: 51.9570,
			validateMinLon: -0.7218,
			validateMaxLon: 0.7218,
			toleranceDeg:   0.02,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLat, maxLat, minLon, maxLon, err := CalculateBoundingBox(tt.lat, tt.lon, tt.radiusKm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if minLat > maxLat {
				t.Errorf("minLat (%f) > maxLat (%f)", minLat, maxLat)
			}

			if minLat < -90 || maxLat > 90 {
				t.Errorf("Latitude out of bounds: minLat=%f, maxLat=%f", minLat, maxLat)
			}

			if minLon < -180 || maxLon > 180 {
				t.Errorf("Longitude out of bounds: minLon=%f, maxLon=%f", minLon, maxLon)
			}

			if !floatEquals(minLat, tt.validateMinLat, tt.toleranceDeg) {
				t.Errorf("minLat = %f, want %f (±%f)", minLat, tt.validateMinLat, tt.toleranceDeg)
			}
			if !floatEquals(maxLat, tt.validateMaxLat, tt.toleranceDeg) {
				t.Errorf("maxLat = %f, want %f (±%f)", maxLat, tt.validateMaxLat, tt.toleranceDeg)
			}
			if !floatEquals(minLon, tt.validateMinLon, tt.toleranceDeg) {
				t.Errorf("minLon = %f, want %f (±%f)", minLon, tt.validateMinLon, tt.toleranceDeg)
			}
			if !floatEquals(maxLon, tt.validateMaxLon, tt.toleranceDeg) {
				t.Errorf("maxLon = %f, want %f (±%f)", maxLon, tt.validateMaxLon, tt.toleranceDeg)
			}
		})
	}
}

func TestCalculateBoundingBoxPoles(t *testing.T) {
	tests := []struct {
		name          string
		lat, lon      float64
		radiusKm      float64
		expectFullLon bool
	}{
		{
			name:          "North Pole",
			lat:           90.0,
			lon:           0.0,
			radiusKm:      100,
			expectFullLon: true,
		},
		{
			name:          "South Pole",
			lat:           -90.0,
			lon:           0.0,
			radiusKm:      100,
			expectFullLon: true,
		},
		{
			name:          "Very close to North Pole",
			lat:           89.9,
			lon:           0.0,
			radiusKm:      50,
			expectFullLon: false,
		},
		{
			name:          "Very close to South Pole",
			lat:           -89.9,
			lon:           0.0,
			radiusKm:      50,
			expectFullLon: false,
		},
		{
			name:          "Large radius from North Pole",
			lat:           90.0,
			lon:           0.0,
			radiusKm:      5000,
			expectFullLon: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLat, maxLat, minLon, maxLon, err := CalculateBoundingBox(tt.lat, tt.lon, tt.radiusKm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if minLat < -90 || minLat > 90 || maxLat < -90 || maxLat > 90 {
				t.Errorf("Latitude out of bounds: [%f, %f]", minLat, maxLat)
			}

			if minLon < -180 || minLon > 180 || maxLon < -180 || maxLon > 180 {
				t.Errorf("Longitude out of bounds: [%f, %f]", minLon, maxLon)
			}

			if tt.expectFullLon {
				if !floatEquals(minLon, -180, floatTolerance) || !floatEquals(maxLon, 180, floatTolerance) {
					t.Errorf("Expected full longitude range [-180, 180], got [%f, %f]", minLon, maxLon)
				}
			}
		})
	}
}

func TestCalculateBoundingBoxDateLine(t *testing.T) {
	tests := []struct {
		name     string
		lat, lon float64
		radiusKm float64
	}{
		{
			name:     "Near date line - positive side",
			lat:      0.0,
			lon:      179.0,
			radiusKm: 100,
		},
		{
			name:     "Near date line - negative side",
			lat:      0.0,
			lon:      -179.0,
			radiusKm: 100,
		},
		{
			name:     "Exactly on date line",
			lat:      0.0,
			lon:      180.0,
			radiusKm: 50,
		},
		{
			name:     "Exactly on negative date line",
			lat:      0.0,
			lon:      -180.0,
			radiusKm: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLat, maxLat, minLon, maxLon, err := CalculateBoundingBox(tt.lat, tt.lon, tt.radiusKm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if minLat >= maxLat {
				t.Errorf("minLat (%f) >= maxLat (%f)", minLat, maxLat)
			}

			if minLat < -90 || maxLat > 90 {
				t.Errorf("Latitude out of bounds: [%f, %f]", minLat, maxLat)
			}

			if minLon < -180 || minLon > 180 || maxLon < -180 || maxLon > 180 {
				t.Errorf("Longitude out of bounds: [%f, %f]", minLon, maxLon)
			}

			t.Logf("Bounding box for (%f, %f) radius %f km: lat[%f, %f], lon[%f, %f]",
				tt.lat, tt.lon, tt.radiusKm, minLat, maxLat, minLon, maxLon)
		})
	}
}

func TestCalculateBoundingBoxLargeRadius(t *testing.T) {
	tests := []struct {
		name     string
		lat, lon float64
		radiusKm float64
	}{
		{
			name:     "Half Earth circumference",
			lat:      0.0,
			lon:      0.0,
			radiusKm: 10000,
		},
		{
			name:     "Quarter Earth from NYC",
			lat:      40.7128,
			lon:      -74.0060,
			radiusKm: 5000,
		},
		{
			name:     "Very large radius",
			lat:      45.0,
			lon:      0.0,
			radiusKm: 15000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLat, maxLat, minLon, maxLon, err := CalculateBoundingBox(tt.lat, tt.lon, tt.radiusKm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if minLat < -90 || maxLat > 90 {
				t.Errorf("Latitude out of bounds: [%f, %f]", minLat, maxLat)
			}
			if minLon < -180 || maxLon > 180 {
				t.Errorf("Longitude out of bounds: [%f, %f]", minLon, maxLon)
			}

			if tt.radiusKm > 10000 {
				latRange := maxLat - minLat
				if latRange < 100 {
					t.Errorf("Expected large latitude range for radius %f km, got %f degrees",
						tt.radiusKm, latRange)
				}
			}
		})
	}
}

func TestCalculateBoundingBoxContainsOriginalPoint(t *testing.T) {
	testCases := []struct {
		lat, lon, radiusKm float64
	}{
		{40.7128, -74.0060, 50},
		{0.0, 0.0, 100},
		{51.5074, -0.1278, 75},
		{-33.8688, 151.2093, 200},
		{35.6762, 139.6503, 150},
	}

	for _, tc := range testCases {
		t.Run("Point within box", func(t *testing.T) {
			minLat, maxLat, minLon, maxLon, err := CalculateBoundingBox(tc.lat, tc.lon, tc.radiusKm)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.lat < minLat || tc.lat > maxLat {
				t.Errorf("Original latitude %f not within bounds [%f, %f]", tc.lat, minLat, maxLat)
			}

			if minLon < maxLon {
				if tc.lon < minLon || tc.lon > maxLon {
					t.Errorf("Original longitude %f not within bounds [%f, %f]", tc.lon, minLon, maxLon)
				}
			}
		})
	}
}

func TestHaversineDistanceValidation(t *testing.T) {
	tests := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
	}{
		{"lat1 too high", 91, 0, 0, 0},
		{"lat1 too low", -91, 0, 0, 0},
		{"lat2 too high", 0, 0, 91, 0},
		{"lat2 too low", 0, 0, -91, 0},
		{"lon1 too high", 0, 181, 0, 0},
		{"lon1 too low", 0, -181, 0, 0},
		{"lon2 too high", 0, 0, 0, 181},
		{"lon2 too low", 0, 0, 0, -181},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestCalculateBoundingBoxValidation(t *testing.T) {
	tests := []struct {
		name     string
		lat, lon float64
		radiusKm float64
	}{
		{"Negative radius", 40.0, -74.0, -10},
		{"Latitude too high", 91.0, 0.0, 50},
		{"Latitude too low", -91.0, 0.0, 50},
		{"Longitude too high", 0.0, 181.0, 50},
		{"Longitude too low", 0.0, -181.0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, err := CalculateBoundingBox(tt.lat, tt.lon, tt.radiusKm)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}
