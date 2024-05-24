package config_test

import (
	"testing"

	"github.com/bitaksi-case/driver-location-api/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("Given an environment name", t, func() {
		env := "test"

		Convey("When loading the config", func() {
			cfg := config.LoadConfig(env)

			Convey("Then the config should be loaded correctly", func() {
				So(cfg.ServerPort, ShouldEqual, "8080")
				So(cfg.ServerHost, ShouldEqual, "localhost")
				So(cfg.SecretKey, ShouldEqual, "valid-token")
				So(cfg.AllowedIssuers, ShouldResemble, []string{"Matching-API", "Driver-API"})
				So(cfg.MongoDBURI, ShouldEqual, "mongodb://localhost:27017")
				So(cfg.MongoDBName, ShouldEqual, "test_db")
				So(cfg.MongoDBCollection, ShouldEqual, "test_collection")
			})
		})
	})
}
