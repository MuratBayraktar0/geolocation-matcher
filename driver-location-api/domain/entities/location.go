package entities

import (
	"math"
)

// Location represents a geographic location with latitude and longitude coordinates.
type Location struct {
	Latitude  float64 // Latitude coordinate
	Longitude float64 // Longitude coordinate
}

// NewLocation creates a new Location with the given latitude and longitude coordinates.
func NewLocation(latitude, longitude float64) *Location {
	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
	}
}

// DistanceTo calculates the great-circle distance between two locations using the haversine formula.
func (loc *Location) DistanceTo(other *Location) float64 {
	// Earth radius in kilometres
	const earthRadius = 6371.0

	// Convert latitude and longitude from degrees to radians
	lat1 := degreesToRadians(loc.Latitude)
	lon1 := degreesToRadians(loc.Longitude)
	lat2 := degreesToRadians(other.Latitude)
	lon2 := degreesToRadians(other.Longitude)

	// Calculate differences in latitude and longitude
	dLat := lat2 - lat1
	dLon := lon2 - lon1

	// Haversine formula for great-circle distance
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate distance
	distance := earthRadius * c

	return distance
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
