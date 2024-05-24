package dtos

import (
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/errors"
)

type DriverCreateRequestDTO struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (dto *DriverCreateRequestDTO) Validate() error {
	if dto.Latitude < -90 || dto.Latitude > 90 {
		return errors.ErrInvalidLatitude
	}
	if dto.Longitude < -180 || dto.Longitude > 180 {
		return errors.ErrInvalidLongitude
	}

	return nil
}

func (dto *DriverCreateRequestDTO) ToLocationEntity() *entities.Driver {
	entity := entities.NewLocation(dto.Latitude, dto.Longitude)
	return entities.NewDriver(dto.ID, entity, time.Now(), time.Now())
}

// DriverBulkCreateRequestDTO swagger:parameters bulkCreateDriverLocation
type DriverBulkCreateRequestDTO struct {
	// in: body
	Drivers []DriverCreateRequestDTO `json:"locations"`
}

func (dto *DriverBulkCreateRequestDTO) Validate() error {
	for _, location := range dto.Drivers {
		if err := location.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (dto *DriverBulkCreateRequestDTO) ToLocationEntities() []*entities.Driver {
	drivers := make([]*entities.Driver, len(dto.Drivers))
	for i, driver := range dto.Drivers {
		drivers[i] = driver.ToLocationEntity()
	}
	return drivers
}

// DriverBulkCreateResponseDTO swagger:response driverLocationsResponse
type DriverBulkCreateResponseDTO struct {
	CreatedIDs []string `json:"created_ids"`
}
