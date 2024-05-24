package dtos

import (
	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/errors"
)

// DriverGetRequestDTO swagger:parameters getDriversLocationbyNear
type DriverGetRequestDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

func NewDriverGetRequestDTO(latitude, longitude, radius float64) *DriverGetRequestDTO {
	return &DriverGetRequestDTO{
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
	}
}

func (dto *DriverGetRequestDTO) Validate() error {
	if dto.Latitude < -90 || dto.Latitude > 90 {
		return errors.ErrInvalidLatitude
	}
	if dto.Longitude < -180 || dto.Longitude > 180 {
		return errors.ErrInvalidLongitude
	}
	if dto.Radius <= 0 {
		return errors.ErrInvalidRadius
	}
	return nil
}

func (dto *DriverGetRequestDTO) ToDriverEntity() *entities.Location {
	return entities.NewLocation(dto.Latitude, dto.Longitude)
}

// DriverGetResponseDTO swagger:response driverLocationsResponse
type DriversGetResponseDTO struct {
	Drivers []DriverGetResponseDTO `json:"drivers"`
}

type DriverGetResponseDTO struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
}
