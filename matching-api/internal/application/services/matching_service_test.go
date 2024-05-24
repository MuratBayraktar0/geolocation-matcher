package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bitaksi-case/matching-api/internal/application/dtos"
	"github.com/bitaksi-case/matching-api/internal/application/services"
	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	. "github.com/smartystreets/goconvey/convey"
)

type mockMatchingService struct {
	matching *entities.Matching
	err      error
}

func (s *mockMatchingService) Matching(ctx context.Context, riderID string, location *entities.Location, radius float64, limit int64) (*entities.Matching, error) {
	return s.matching, s.err
}

func TestMatchingService_Match(t *testing.T) {
	Convey("Given a MatchingService", t, func() {
		mockService := &mockMatchingService{
			matching: &entities.Matching{},
			err:      nil,
		}
		matchingService := services.NewMatchingService(mockService)

		Convey("When Match is called with a valid DTO", func() {
			riderID := "rider123"
			dto := &dtos.MatchingRequestDTO{
				RiderID:   riderID,
				Latitude:  40.7128,
				Longitude: -74.0060,
			}
			radius := 10.0
			limit := int64(5)

			Convey("And the service returns a non-empty matching", func() {
				mockService.matching = &entities.Matching{
					RiderID: riderID,
					Matching: []*entities.Match{
						{
							DriverID: "driver123",
							Location: &entities.Location{
								Latitude:  40.7128,
								Longitude: -74.0060,
							},
							Distance: 5.0,
						},
					},
				}

				response, err := matchingService.Match(context.Background(), dto, radius, limit)

				Convey("Then the response should match the expected DTO", func() {
					So(err, ShouldBeNil)
					So(response, ShouldNotBeNil)
					So(response.RiderID, ShouldEqual, riderID)
					So(response.Matching[0].DriverID, ShouldEqual, "driver123")
					So(response.Matching[0].Distance, ShouldEqual, 5.0)
					So(response.Matching[0].Latitude, ShouldEqual, 40.7128)
					So(response.Matching[0].Longitude, ShouldEqual, -74.0060)
				})
			})

			Convey("And the service returns an empty matching", func() {
				mockService.matching = &entities.Matching{
					RiderID:  riderID,
					Matching: []*entities.Match{},
				}

				Convey("And the service returns no error", func() {
					response, err := matchingService.Match(context.Background(), dto, radius, limit)

					Convey("Then the response should match the expected DTO", func() {
						So(err, ShouldBeNil)
						So(response, ShouldNotBeNil)
					})
				})

				Convey("And the service returns an error", func() {
					someError := errors.New("some error")
					mockService.err = someError

					response, err := matchingService.Match(context.Background(), dto, radius, limit)

					Convey("Then the response should be nil and the error should match", func() {
						So(err, ShouldEqual, someError)
						So(response, ShouldBeNil)
					})
				})
			})
		})
	})
}
