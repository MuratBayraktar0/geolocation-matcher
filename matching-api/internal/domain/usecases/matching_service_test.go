package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/domain/usecases"
	. "github.com/smartystreets/goconvey/convey"
)

type mockDriverLocationService struct {
	drivers *entities.Drivers
}

func (s *mockDriverLocationService) GetDriversLocationbyNear(ctx context.Context, point *entities.Location, radius float64, limit int64) (*entities.Drivers, error) {
	return s.drivers, nil
}

type mockMatchingRepository struct {
	matching *entities.Matching
}

func (r *mockMatchingRepository) GetMatching(ctx context.Context, riderID string, latitude, longitude, radius float64, limit int) (*entities.Matching, error) {
	return r.matching, nil
}

func (r *mockMatchingRepository) UpdateDriverLocation(ctx context.Context, driverID string, latitude, longitude float64) error {
	return nil
}

func TestMatchingService_Matching(t *testing.T) {
	Convey("Given a MatchingService", t, func() {
		mockDriverLocationService := &mockDriverLocationService{
			drivers: &entities.Drivers{
				Drivers: []*entities.Driver{
					{
						ID:        "driver1234",
						Latitude:  40.7128,
						Longitude: -74.0060,
						Distance:  5.0,
					},
				},
			},
		}
		mockMatchingRepository := &mockMatchingRepository{
			matching: &entities.Matching{},
		}
		matchingService := usecases.NewMatchingService(mockDriverLocationService, mockMatchingRepository)

		Convey("When Matching is called", func() {
			riderID := "rider123"
			location := &entities.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			}
			radius := 10.0
			limit := int64(5)

			Convey("And GetMatching returns a non-empty cacheMatching", func() {
				mockMatchingRepository.matching.Matching = []*entities.Match{
					{
						DriverID: "driver123",
						Location: &entities.Location{
							Latitude:  40.7128,
							Longitude: -74.0060,
						},
						Distance: 5.0,
					},
				}
				mockMatchingRepository.matching.RiderID = riderID

				matching, err := matchingService.Matching(context.Background(), riderID, location, radius, limit)
				So(err, ShouldBeNil)
				So(matching, ShouldEqual, mockMatchingRepository.matching)
			})

			Convey("And GetMatching returns an empty cacheMatching", func() {
				mockMatchingRepository.matching.Matching = []*entities.Match{}
				Convey("And GetDriversLocationbyNear returns a list of drivers", func() {
					matching, err := matchingService.Matching(context.Background(), riderID, location, radius, limit)
					fmt.Println(matching)
					So(err, ShouldBeNil)
					So(matching.RiderID, ShouldEqual, riderID)
					So(len(matching.Matching), ShouldEqual, 1)
					So(matching.Matching[0].DriverID, ShouldEqual, "driver1234")
				})

			})
		})
	})
}
