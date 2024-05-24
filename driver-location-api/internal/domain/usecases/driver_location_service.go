package usecases

import (
	"context"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/interfaces"
)

type driverLocationService struct {
	repo interfaces.DriverLocationRepository
}

func NewDriverLocationService(repo interfaces.DriverLocationRepository) interfaces.DriverLocationService {
	return &driverLocationService{repo: repo}
}

func (s *driverLocationService) BulkCreateDriverLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error) {
	result, err := s.repo.BulkCreateLocation(ctx, drivers)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *driverLocationService) GetDriversLocationbyNear(ctx context.Context, latitude, longitude, radius float64, limit int64) (*entities.Drivers, error) {
	location := entities.NewLocation(latitude, longitude)
	drivers, err := s.repo.GetDriversLocationbyNear(ctx, location, radius, limit)
	if err != nil {
		return nil, err
	}

	return drivers, nil
}
