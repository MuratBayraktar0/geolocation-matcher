package config_test

import (
	"testing"

	"github.com/bitaksi-case/matching-api/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("Given an environment name", t, func() {
		env := "test"

		Convey("When loading the config", func() {
			cfg := config.LoadConfig(env)

			Convey("Then the config should be loaded correctly", func() {
				So(cfg.ServerPort, ShouldEqual, "8081")
				So(cfg.ServerHost, ShouldEqual, "localhost")
				So(cfg.SecretKey, ShouldEqual, "valid-token")
				So(cfg.AllowedIssuers, ShouldResemble, []string{"User-API"})
				So(cfg.RedisAddr, ShouldEqual, "localhost:6379")
				So(cfg.RedisPassword, ShouldEqual, "")
				So(cfg.RedisDB, ShouldEqual, 0)
				So(cfg.AuthAPIEndpoint, ShouldEqual, "http://localhost:8082")
				So(cfg.DriverLocationApiEndpoint, ShouldEqual, "http://localhost:8080")
			})
		})
	})
}
