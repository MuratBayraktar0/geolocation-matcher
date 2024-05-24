package services

import (
	"context"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/application/dtos"
	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationService struct {
	service interfaces.DriverLocationService
}

func NewLocationService(service interfaces.DriverLocationService) *LocationService {
	return &LocationService{
		service: service,
	}
}

func (s *LocationService) BulkCreateDriverLocation(ctx context.Context, dto *dtos.DriverBulkCreateRequestDTO) (*dtos.DriverBulkCreateResponseDTO, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	drivers := dto.ToLocationEntities()
	for _, driver := range drivers {
		if driver.ID == "" {
			driver.ID = primitive.NewObjectID().Hex()
			driver.CreateAt = time.Now()
			driver.UpdateAt = time.Now()
		} else {
			driver.UpdateAt = time.Now()
		}
	}

	IDs, err := s.service.BulkCreateDriverLocation(ctx, entities.NewDrivers(drivers))
	if err != nil {
		return nil, err
	}
	return &dtos.DriverBulkCreateResponseDTO{CreatedIDs: *IDs}, nil
}

// GetLocation retrieves a location by its ID.
func (s *LocationService) GetDriversLocationbyNear(ctx context.Context, dto *dtos.DriverGetRequestDTO, limit int64) (*dtos.DriversGetResponseDTO, error) {
	if limit <= 0 {
		limit = 10
	}

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	locations, err := s.service.GetDriversLocationbyNear(ctx, dto.Latitude, dto.Longitude, dto.Radius, limit)
	if err != nil {
		return nil, err
	}

	return ToLocationsResponseDTO(locations), nil
}

// ToLocationsResponseDTO converts a Locations entity to a LocationsResponseDTO.
func ToLocationsResponseDTO(locationsEntity *entities.Drivers) *dtos.DriversGetResponseDTO {
	locations := []dtos.DriverGetResponseDTO{}
	for _, driver := range locationsEntity.Drivers {
		locations = append(locations, dtos.DriverGetResponseDTO{
			ID:        driver.ID,
			Latitude:  driver.Location.Coordinates[0],
			Longitude: driver.Location.Coordinates[1],
			Distance:  driver.Distance,
		})
	}

	return &dtos.DriversGetResponseDTO{Drivers: locations}
}
