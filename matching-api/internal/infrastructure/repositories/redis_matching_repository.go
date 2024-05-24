package repositories

import (
	"context"
	"time"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/domain/interfaces"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/go-redis/redis/v8"
)

type RedisDriverLocationRepository struct {
	adapter *adapters.RedisClient
}

func NewRedisDriverLocationRepository(adapter *adapters.RedisClient) interfaces.MatchingRepository {
	return &RedisDriverLocationRepository{adapter: adapter}
}

func (repo *RedisDriverLocationRepository) UpdateDriverLocation(ctx context.Context, driverID string, latitude, longitude float64) error {
	err := repo.adapter.Client.GeoAdd(ctx, "drivers", &redis.GeoLocation{
		Name:      driverID,
		Longitude: longitude,
		Latitude:  latitude,
	}).Err()
	if err != nil {
		return err
	}

	return repo.adapter.Client.Expire(ctx, "drivers", 10*time.Second).Err()
}

func (repo *RedisDriverLocationRepository) GetMatching(ctx context.Context, riderID string, latitude, longitude, radius float64, limit int) (*entities.Matching, error) {
	locations, err := repo.adapter.Client.GeoRadius(ctx, "drivers", longitude, latitude, &redis.GeoRadiusQuery{
		Radius:   radius,
		Unit:     "km",
		WithDist: true,
		Sort:     "ASC",
		Count:    limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	matches := []*entities.Match{}
	for _, location := range locations {
		match := entities.NewMatch(location.Name, entities.NewLocation(location.Latitude, location.Longitude), location.Dist, time.Now(), time.Now())
		matches = append(matches, match)
	}
	return entities.NewMatching(riderID, matches), nil
}
