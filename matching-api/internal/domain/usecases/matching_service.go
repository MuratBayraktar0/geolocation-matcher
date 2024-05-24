package usecases

import (
	"context"
	"time"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/domain/interfaces"
)

type matchingService struct {
	service interfaces.DriverLocationService
	repo    interfaces.MatchingRepository
}

func NewMatchingService(service interfaces.DriverLocationService, repo interfaces.MatchingRepository) interfaces.MatchingService {
	return &matchingService{service: service, repo: repo}
}

func (s *matchingService) Matching(ctx context.Context, riderID string, location *entities.Location, radius float64, limit int64) (*entities.Matching, error) {
	cacheMatching, err := s.repo.GetMatching(ctx, riderID, location.Latitude, location.Longitude, radius, int(limit))
	if err != nil {
		return nil, err
	}

	if len(cacheMatching.Matching) > 0 {
		return cacheMatching, nil
	}

	destination := entities.NewLocation(location.Latitude, location.Longitude)
	drivers, err := s.service.GetDriversLocationbyNear(ctx, destination, radius, limit)
	if err != nil {
		return nil, err
	}

	matching := []*entities.Match{}
	for _, driver := range drivers.Drivers {
		origin := entities.NewLocation(driver.Latitude, driver.Longitude)
		match := entities.NewMatch(driver.ID, origin, driver.Distance, time.Now(), time.Now())
		err := s.repo.UpdateDriverLocation(ctx, driver.ID, driver.Latitude, driver.Longitude)
		if err != nil {
			return nil, err
		}

		matching = append(matching, match)
	}

	return entities.NewMatching(riderID, matching), nil
}
