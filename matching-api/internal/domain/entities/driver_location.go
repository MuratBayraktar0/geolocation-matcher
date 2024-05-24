package entities

type Drivers struct {
	Drivers []*Driver
}

type Driver struct {
	ID        string
	Latitude  float64
	Longitude float64
	Distance  float64
}
