package geo

import (
	"fmt"
	"math"
)

const (
	EarthRadiusKm = 6371.0

	minLatRad = -math.Pi / 2
	maxLatRad = math.Pi / 2
	minLonRad = -math.Pi
	maxLonRad = math.Pi
)

// CalculateBoundingBox calculates the bounding box for a circle defined by a center point
// (lat, lon in degrees) and a radius in kilometers.
// It returns the minimum and maximum latitudes and longitudes that define the bounding box.
func CalculateBoundingBox(lat, lon, radiusKm float64) (minLat, maxLat, minLon, maxLon float64, err error) {
	if radiusKm < 0 {
		return 0, 0, 0, 0, fmt.Errorf("radius cannot be negative")
	}
	if lat < -90 || lat > 90 {
		return 0, 0, 0, 0, fmt.Errorf("latitude must be between -90 and 90 degrees")
	}
	if lon < -180 || lon > 180 {
		return 0, 0, 0, 0, fmt.Errorf("longitude must be between -180 and 180 degrees")
	}

	latRad := degToRad(lat)
	lonRad := degToRad(lon)
	angularDistance := radiusKm / EarthRadiusKm

	minLatResult := latRad - angularDistance
	maxLatResult := latRad + angularDistance

	if minLatResult > minLatRad && maxLatResult < maxLatRad {
		deltaLon := math.Asin(math.Sin(angularDistance) / math.Cos(latRad))

		minLonResult := lonRad - deltaLon
		maxLonResult := lonRad + deltaLon

		if minLonResult < minLonRad {
			minLonResult += 2 * math.Pi
		}
		if maxLonResult > maxLonRad {
			maxLonResult -= 2 * math.Pi
		}

		return radToDeg(minLatResult), radToDeg(maxLatResult), radToDeg(minLonResult), radToDeg(maxLonResult), nil
	}

	minLatResult = math.Max(minLatResult, minLatRad)
	maxLatResult = math.Min(maxLatResult, maxLatRad)

	return radToDeg(minLatResult), radToDeg(maxLatResult), radToDeg(minLonRad), radToDeg(maxLonRad), nil
}

// HaversineDistance calculates the great-circle distance between two points
// (lat1, lon1) and (lat2, lon2) in kilometers using the Haversine formula.
// Source: https://www.movable-type.co.uk/scripts/latlong.html
func HaversineDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	if lat1 < -90 || lat1 > 90 || lat2 < -90 || lat2 > 90 {
		return 0, fmt.Errorf("latitude must be between -90 and 90 degrees")
	}
	if lon1 < -180 || lon1 > 180 || lon2 < -180 || lon2 > 180 {
		return 0, fmt.Errorf("longitude must be between -180 and 180 degrees")
	}

	lat1Rad := degToRad(lat1)
	lon1Rad := degToRad(lon1)
	lat2Rad := degToRad(lat2)
	lon2Rad := degToRad(lon2)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	sinHalfDeltaLat := math.Sin(deltaLat / 2)
	sinHalfDeltaLon := math.Sin(deltaLon / 2)
	a := sinHalfDeltaLat*sinHalfDeltaLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinHalfDeltaLon*sinHalfDeltaLon
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c, nil
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}
