package entities_test

import (
	"testing"

	"github.com/bitaksi-case/driver-location-api/domain/entities"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocation(t *testing.T) {
	Convey("Given a location", t, func() {
		loc1 := entities.NewLocation(40.7128, -74.0060)
		loc2 := entities.NewLocation(34.0522, -118.2437)

		Convey("When creating a new location", func() {
			Convey("The location should be created with the correct coordinates", func() {
				So(loc1, ShouldNotBeNil)
				So(loc1.Latitude, ShouldEqual, 40.7128)
				So(loc1.Longitude, ShouldEqual, -74.0060)
			})
		})

		Convey("When calculating distance to another location", func() {
			distance := loc1.DistanceTo(loc2)
			expectedDistance := 3935.7 // kilometres

			Convey("Should return the expected distance", func() {
				So(distance, ShouldAlmostEqual, expectedDistance, 0.1)
			})
		})
	})
}
