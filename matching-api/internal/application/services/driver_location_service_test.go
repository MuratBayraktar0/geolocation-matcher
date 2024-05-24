package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/bitaksi-case/matching-api/internal/application/services"
	"github.com/bitaksi-case/matching-api/internal/config"
	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
)

func TestDriverLocationService_GetDriversLocationbyNear(t *testing.T) {
	Convey("Given a DriverLocationService", t, func() {
		ctx := context.Background()
		cfg := config.LoadConfig("test")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"status": 200,
				"data": map[string]interface{}{
					"drivers": []map[string]interface{}{
						{
							"id":        "664bd3cc4384b4cbc5551a50",
							"latitude":  41.06040004,
							"longitude": 27.83085102,
							"distance":  11336386.96747813,
						},
						{
							"id":        "664bd3cc4384b4cbc55521bc",
							"latitude":  41.07653184,
							"longitude": 27.83089933,
							"distance":  11336395.120692058,
						},
					},
				},
				"error": nil,
				"meta": map[string]string{
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				},
			}

			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
		}))
		defer ts.Close()

		cfg.DriverLocationApiEndpoint = ts.URL

		httpClientAdapter := adapters.NewHttpClient(ctx)
		authClientAdapter := adapters.NewAuthClient(cfg.AuthAPIEndpoint)

		service := services.NewDriverLocationService(httpClientAdapter, authClientAdapter, cfg)

		Convey("When calling GetDriversLocationbyNear", func() {
			point := &entities.Location{
				Latitude:  40.7128,
				Longitude: -74.0060,
			}
			radius := 10.0
			limit := int64(5)
			var drivers *entities.Drivers
			drivers, err := service.GetDriversLocationbyNear(context.Background(), point, radius, limit)

			Convey("Then the drivers should be retrieved successfully", func() {
				So(err, ShouldBeNil)
				So(drivers, ShouldNotBeNil)
				So(len(drivers.Drivers), ShouldEqual, 2)
				So(drivers.Drivers[0].ID, ShouldEqual, "664bd3cc4384b4cbc5551a50")
				So(drivers.Drivers[0].Latitude, ShouldEqual, 41.06040004)
				So(drivers.Drivers[0].Longitude, ShouldEqual, 27.83085102)
				So(drivers.Drivers[0].Distance, ShouldEqual, 11336386.96747813)
				So(drivers.Drivers[1].ID, ShouldEqual, "664bd3cc4384b4cbc55521bc")
				So(drivers.Drivers[1].Latitude, ShouldEqual, 41.07653184)
				So(drivers.Drivers[1].Longitude, ShouldEqual, 27.83089933)
				So(drivers.Drivers[1].Distance, ShouldEqual, 11336395.120692058)
			})
		})
	})
}
