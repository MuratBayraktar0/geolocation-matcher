package interfaces

import (
	"context"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
)

// MatchingRepository defines the methods that any
// data storage provider needs to implement to get and save driver locations.
type MatchingRepository interface {
	GetMatching(ctx context.Context, riderID string, latitude, longitude, radius float64, limit int) (*entities.Matching, error)
	UpdateDriverLocation(ctx context.Context, driverID string, latitude, longitude float64) error
}
