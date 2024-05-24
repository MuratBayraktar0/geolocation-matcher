package interfaces

import (
	"context"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
)

// DriverLocationRepository defines the methods that any
// data storage provider needs to implement to get and save driver locations.
type DriverLocationRepository interface {
	BulkCreateLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error)
	GetDriversLocationbyNear(ctx context.Context, location *entities.Location, radius float64, limit int64) (*entities.Drivers, error)
	GetLocationCount(ctx context.Context) (int64, error)
	CreateLocationIndex(ctx context.Context) error
}
