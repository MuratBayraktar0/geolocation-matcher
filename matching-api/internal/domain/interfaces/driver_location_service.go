package interfaces

import (
	"context"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
)

// DriverLocationService defines the service for managing driver locations.
type DriverLocationService interface {
	GetDriversLocationbyNear(ctx context.Context, point *entities.Location, radius float64, limit int64) (*entities.Drivers, error)
}
