package entities

import (
	"time"
)

type Matching struct {
	RiderID  string   `json:"rider_id"`
	Matching []*Match `json:"matching"`
}

func NewMatching(riderID string, matching []*Match) *Matching {
	return &Matching{
		RiderID:  riderID,
		Matching: matching,
	}
}

type Match struct {
	DriverID  string    `json:"driver_id"`
	Location  *Location `json:"location"`
	Distance  float64   `json:"distance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewMatch(driverID string, location *Location, distance float64, createdAt, updatedAt time.Time) *Match {
	return &Match{
		DriverID:  driverID,
		Location:  location,
		Distance:  distance,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLocation(latitude, longitude float64) *Location {
	return &Location{
		Latitude:  latitude,
		Longitude: longitude,
	}
}
