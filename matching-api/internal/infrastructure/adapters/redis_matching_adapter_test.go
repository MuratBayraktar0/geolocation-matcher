package adapters

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRedisClient(t *testing.T) {
	Convey("Given Redis client parameters", t, func() {
		addr := "localhost:6379"
		password := "password"
		db := 0

		Convey("When creating a new Redis client", func() {
			redisClient := NewRedisClient(addr, password, db)

			Convey("Then the Redis client should be created successfully", func() {
				So(redisClient, ShouldNotBeNil)
				So(redisClient.Client, ShouldNotBeNil)
			})
		})
	})
}
