package interfaces

import (
	"context"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
)

// MatchingService defines the service for matching riders with drivers.
type MatchingService interface {
	Matching(ctx context.Context, riderID string, location *entities.Location, radius float64, limit int64) (*entities.Matching, error)
}
