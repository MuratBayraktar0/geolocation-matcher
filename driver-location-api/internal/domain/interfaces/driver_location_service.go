package interfaces

import (
	"context"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
)

type DriverLocationService interface {
	BulkCreateDriverLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error)
	GetDriversLocationbyNear(ctx context.Context, latitude, longitude, radius float64, limit int64) (*entities.Drivers, error)
}
