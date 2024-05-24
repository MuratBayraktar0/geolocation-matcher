package repositories_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bitaksi-case/driver-location-api/internal/config"
	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/repositories"
)

func TestMongoDBDriverLocationRepository(t *testing.T) {
	Convey("Given a MongoDBDriverLocationRepository", t, func() {
		ctx := context.Background()
		cfg := config.LoadConfig("test")
		client, _ := adapters.NewMongoDBClient(ctx, cfg.MongoDBURI)
		dbName := "test_db"
		colName := "test_collection"
		repo := repositories.NewMongoDBDriverLocationRepository(client, dbName, colName)

		Convey("When BulkCreateLocation is called", func() {
			ctx := context.Background()
			drivers := &entities.Drivers{
				Drivers: []*entities.Driver{
					{
						ID:       primitive.NewObjectID().Hex(),
						Location: entities.NewLocation(40.7128, -74.0060),
						Distance: 0,
						UpdateAt: time.Now(),
						CreateAt: time.Now(),
					},
					{
						ID:       primitive.NewObjectID().Hex(),
						Location: entities.NewLocation(-118.244, 34.0522),
						Distance: 0,
						UpdateAt: time.Now(),
						CreateAt: time.Now(),
					},
				},
			}

			insertedIDs, err := repo.BulkCreateLocation(ctx, drivers)

			Convey("Then the drivers should be inserted successfully", func() {
				So(err, ShouldBeNil)
				So(insertedIDs, ShouldNotBeNil)
				So(len(*insertedIDs), ShouldEqual, 2)
			})
		})

		Convey("When GetDriversLocationbyNear is called", func() {
			ctx := context.Background()
			loc := entities.NewLocation(-118.244, 34.0522)
			radius := 10.0
			limit := int64(5)
			err := repo.CreateLocationIndex(ctx)
			So(err, ShouldBeNil)
			drivers, err := repo.GetDriversLocationbyNear(ctx, loc, radius, limit)

			Convey("Then the drivers within the specified radius should be returned", func() {
				So(err, ShouldBeNil)
				So(drivers, ShouldNotBeNil)
				So(len(drivers.Drivers), ShouldBeGreaterThan, 0)
			})
		})

		Convey("When GetLocationCount is called", func() {
			ctx := context.Background()

			count, err := repo.GetLocationCount(ctx)

			Convey("Then the count of locations should be returned", func() {
				So(err, ShouldBeNil)
				So(count, ShouldBeGreaterThanOrEqualTo, 0)
			})
		})

		Convey("When CreateLocationIndex is called", func() {
			ctx := context.Background()

			err := repo.CreateLocationIndex(ctx)

			Convey("Then the location index should be created successfully", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
