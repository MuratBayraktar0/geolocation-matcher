package entities

import (
	"time"
)

type Drivers struct {
	Drivers []*Driver
}

func NewDrivers(drivers []*Driver) *Drivers {
	return &Drivers{
		Drivers: drivers,
	}
}

type Driver struct {
	ID       string    `bson:"_id,omitempty"`
	Location *Location `bson:"location"`
	Distance float64   `bson:"distance,omitempty"`
	UpdateAt time.Time `bson:"update_at"`
	CreateAt time.Time `bson:"create_at"`
}

func NewDriver(id string, location *Location, updateAt, createAt time.Time) *Driver {
	return &Driver{
		ID:       id,
		Location: location,
		UpdateAt: updateAt,
		CreateAt: createAt,
	}
}

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

func NewLocation(latitude, longitude float64) *Location {
	return &Location{
		Type:        "Point",
		Coordinates: []float64{latitude, longitude},
	}
}
