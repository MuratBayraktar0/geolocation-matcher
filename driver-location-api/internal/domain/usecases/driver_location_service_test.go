package usecases_test

import (
	"context"
	"testing"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/usecases"
	. "github.com/smartystreets/goconvey/convey"
)

type mockDriverLocationRepo struct {
	drivers *entities.Drivers
}

func (m *mockDriverLocationRepo) CreateLocationIndex(ctx context.Context) error {
	return nil
}

func (m *mockDriverLocationRepo) BulkCreateLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error) {
	return nil, nil
}

func (m *mockDriverLocationRepo) GetDriversLocationbyNear(ctx context.Context, location *entities.Location, radius float64, limit int64) (*entities.Drivers, error) {
	return m.drivers, nil
}

func (m *mockDriverLocationRepo) GetLocationCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func TestDriverLocationService_GetDriversLocationbyNear(t *testing.T) {
	Convey("Given a DriverLocationService", t, func() {
		mockRepo := &mockDriverLocationRepo{
			drivers: &entities.Drivers{
				Drivers: []*entities.Driver{
					{
						ID:       "driver1",
						Location: entities.NewLocation(40.7128, -74.0060),
						Distance: 5.0,
					},
				},
			},
		}

		service := usecases.NewDriverLocationService(mockRepo)

		Convey("When GetDriversLocationbyNear is called", func() {
			ctx := context.TODO()
			point := entities.NewLocation(40.7128, -74.0060)
			radius := 10.0
			limit := int64(10)

			Convey("Then it should return a list of drivers", func() {
				drivers, err := service.GetDriversLocationbyNear(ctx, point.Coordinates[0], point.Coordinates[1], radius, limit)
				So(err, ShouldBeNil)
				So(drivers.Drivers[0].ID, ShouldEqual, "driver1")
			})
		})
	})
}
