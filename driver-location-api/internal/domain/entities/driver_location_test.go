package entities_test

import (
	"testing"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLocation(t *testing.T) {
	Convey("Given a location", t, func() {
		loc1 := entities.NewLocation(40.7128, -74.0060)

		Convey("When creating a new location", func() {
			Convey("The location should be created with the correct coordinates", func() {
				So(loc1, ShouldNotBeNil)
				So(loc1.Type, ShouldEqual, "Point")
				So(loc1.Coordinates[0], ShouldEqual, 40.7128)
				So(loc1.Coordinates[1], ShouldEqual, -74.0060)
			})
		})
	})
}

func TestDrivers(t *testing.T) {
	Convey("Given a list of drivers", t, func() {
		driver1 := entities.NewDriver("1", entities.NewLocation(40.7128, -74.0060), time.Now(), time.Now())
		driver2 := entities.NewDriver("2", entities.NewLocation(37.7749, -122.4194), time.Now(), time.Now())
		drivers := []*entities.Driver{driver1, driver2}
		driversObj := entities.NewDrivers(drivers)

		Convey("When creating a new drivers object", func() {
			Convey("The drivers object should be created with the correct drivers", func() {
				So(driversObj, ShouldNotBeNil)
				So(driversObj.Drivers, ShouldHaveLength, 2)
				So(driversObj.Drivers[0].ID, ShouldEqual, "1")
				So(driversObj.Drivers[0].Location.Type, ShouldEqual, "Point")
				So(driversObj.Drivers[0].Location.Coordinates[0], ShouldEqual, 40.7128)
				So(driversObj.Drivers[0].Location.Coordinates[1], ShouldEqual, -74.0060)
				So(driversObj.Drivers[1].ID, ShouldEqual, "2")
				So(driversObj.Drivers[1].Location.Type, ShouldEqual, "Point")
				So(driversObj.Drivers[1].Location.Coordinates[0], ShouldEqual, 37.7749)
				So(driversObj.Drivers[1].Location.Coordinates[1], ShouldEqual, -122.4194)
			})
		})
	})
}
