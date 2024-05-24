package entities_test

import (
	"testing"
	"time"

	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMatching(t *testing.T) {
	Convey("Given a new matching", t, func() {
		// Create test data
		riderID := "1"
		origin := entities.NewLocation(40.7128, -74.0060)
		distance := 100.0
		createdAt := time.Now()
		updatedAt := time.Now()

		// Create a match
		match := entities.NewMatch("1", origin, distance, createdAt, updatedAt)

		// Call the function
		matching := entities.NewMatching(riderID, []*entities.Match{match})

		Convey("The matching should have the correct values", func() {
			So(matching.RiderID, ShouldEqual, riderID)
			So(matching.Matching, ShouldResemble, []*entities.Match{match})
			So(matching.Matching[0].DriverID, ShouldEqual, "1")
			So(matching.Matching[0].Location, ShouldEqual, origin)
			So(matching.Matching[0].Distance, ShouldEqual, distance)
			So(matching.Matching[0].CreatedAt, ShouldEqual, createdAt)
			So(matching.Matching[0].UpdatedAt, ShouldEqual, updatedAt)
		})

		Convey("The match should have the correct values", func() {
			So(match.DriverID, ShouldEqual, "1")
			So(match.Location, ShouldEqual, origin)
			So(match.Distance, ShouldEqual, distance)
			So(match.CreatedAt, ShouldEqual, createdAt)
			So(match.UpdatedAt, ShouldEqual, updatedAt)
		})

		Convey("The location should have the correct values", func() {
			So(origin.Latitude, ShouldEqual, 40.7128)
			So(origin.Longitude, ShouldEqual, -74.0060)
		})
	})
}
