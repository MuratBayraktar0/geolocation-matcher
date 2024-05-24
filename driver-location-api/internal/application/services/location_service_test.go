package services_test

import (
	"context"
	"testing"

	"github.com/bitaksi-case/driver-location-api/internal/application/dtos"
	"github.com/bitaksi-case/driver-location-api/internal/application/services"
	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	. "github.com/smartystreets/goconvey/convey"
)

type mockDriverLocationService struct {
	drivers *entities.Drivers
}

func (m *mockDriverLocationService) BulkCreateDriverLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error) {
	driversList := make([]string, 0)
	for _, driver := range drivers.Drivers {
		driversList = append(driversList, driver.ID)
	}
	return &driversList, nil
}

func (m *mockDriverLocationService) GetDriversLocationbyNear(ctx context.Context, latitude, longitude, radius float64, limit int64) (*entities.Drivers, error) {
	return m.drivers, nil
}

func TestLocationService(t *testing.T) {
	Convey("Given a LocationService", t, func() {
		mockDriverLocationService := &mockDriverLocationService{
			drivers: &entities.Drivers{
				Drivers: []*entities.Driver{
					{
						ID:       "driver1",
						Location: entities.NewLocation(40.7128, -74.0060),
						Distance: 5.0,
					},
					{
						ID:       "driver2",
						Location: entities.NewLocation(34.0522, -118.2437),
						Distance: 10.0,
					},
				},
			},
		}
		locationService := services.NewLocationService(mockDriverLocationService)

		Convey("When BulkCreateDriverLocation is called with valid input", func() {
			dto := &dtos.DriverBulkCreateRequestDTO{
				Drivers: []dtos.DriverCreateRequestDTO{
					{
						ID:        "id1",
						Latitude:  40.7128,
						Longitude: -74.0060,
					},
					{
						ID:        "id2",
						Latitude:  34.0522,
						Longitude: -118.2437,
					},
				},
			}
			expectedIDs := []string{"id1", "id2"}

			Convey("It should create driver locations in bulk and return the created IDs", func() {
				createdIDs, err := locationService.BulkCreateDriverLocation(context.Background(), dto)
				So(err, ShouldBeNil)
				So(createdIDs.CreatedIDs, ShouldResemble, expectedIDs)
			})
		})

		Convey("When GetDriversLocationbyNear is called with valid input", func() {
			limit := int64(10)
			radius := 10.0
			dto := &dtos.DriverGetRequestDTO{
				Latitude:  40.7128,
				Longitude: -74.0060,
				Radius:    radius,
			}
			expectedLocations := &entities.Drivers{
				Drivers: []*entities.Driver{
					{
						ID:       "driver1",
						Location: entities.NewLocation(40.7128, -74.0060),
						Distance: 5.0,
					},
				},
			}

			Convey("It should retrieve drivers' locations near the specified coordinates", func() {
				locations, err := locationService.GetDriversLocationbyNear(context.Background(), dto, limit)
				So(err, ShouldBeNil)
				So(locations.Drivers[0].ID, ShouldEqual, expectedLocations.Drivers[0].ID)
			})
		})
	})
}
