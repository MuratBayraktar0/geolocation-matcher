package repositories_test

import (
	"context"
	"testing"

	"github.com/bitaksi-case/matching-api/internal/config"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/repositories"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisDriverLocationRepository_UpdateDriverLocation(t *testing.T) {
	Convey("Given a Redis driver location repository", t, func() {
		cfg := config.LoadConfig("test")
		adapter := adapters.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
		repo := repositories.NewRedisDriverLocationRepository(adapter)

		Convey("When updating driver location", func() {
			ctx := context.TODO()
			driverID := "driver-123"
			latitude := 40.7128
			longitude := -74.0060

			err := repo.UpdateDriverLocation(ctx, driverID, latitude, longitude)

			Convey("Then the driver location should be updated successfully", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestRedisDriverLocationRepository_GetMatching(t *testing.T) {
	Convey("Given a Redis driver location repository", t, func() {
		cfg := config.LoadConfig("test")
		adapter := adapters.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
		repo := repositories.NewRedisDriverLocationRepository(adapter)

		Convey("When getting matching drivers", func() {
			ctx := context.TODO()
			riderID := "rider-123"
			latitude := 40.7128
			longitude := -74.0060
			radius := 10.0
			limit := 5

			matching, err := repo.GetMatching(ctx, riderID, latitude, longitude, radius, limit)

			Convey("Then the matching drivers should be retrieved successfully", func() {
				So(err, ShouldBeNil)
				So(matching, ShouldNotBeNil)
				So(matching.RiderID, ShouldEqual, riderID)
				So(len(matching.Matching), ShouldBeGreaterThan, 0)
			})
		})
	})
}
