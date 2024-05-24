package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/application/dtos"
	"github.com/bitaksi-case/driver-location-api/internal/application/handlers"
	"github.com/bitaksi-case/driver-location-api/internal/application/services"
	"github.com/bitaksi-case/driver-location-api/internal/config"
	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/usecases"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/repositories"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/server"
	"github.com/dgrijalva/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
)

type Respose struct {
	Status int
	Data   interface{}
	Error  string
	Meta   map[string]string
}

func TestIntegration(t *testing.T) {
	Convey("Given a new HTTP server", t, func() {
		ctx := context.Background()
		cfg := config.LoadConfig("test")

		mongoAdapter, _ := adapters.NewMongoDBClient(ctx, cfg.MongoDBURI)
		repo := repositories.NewMongoDBDriverLocationRepository(mongoAdapter, cfg.MongoDBName, cfg.MongoDBCollection)
		repo.BulkCreateLocation(ctx, &entities.Drivers{
			Drivers: []*entities.Driver{
				{
					ID:       "1",
					Location: entities.NewLocation(40.7128, -74.0060),
				}},
		})
		coreService := usecases.NewDriverLocationService(repo)
		service := services.NewLocationService(coreService)
		handler := handlers.NewDriverLocationHandler(ctx, service)
		server := server.NewHTTPServer(cfg)

		Convey("When registering a handler for POST /driver/locations", func() {
			app := server.RegisterHandler(http.MethodPost, "/driver/locations", handler.BulkCreateDriverLocation)
			dto := dtos.DriverBulkCreateRequestDTO{
				Drivers: []dtos.DriverCreateRequestDTO{
					{
						Latitude:  41.06040004,
						Longitude: 27.83085102,
					},
					{
						Latitude:  41.06040004,
						Longitude: 27.83085102,
					},
				},
			}
			body, _ := json.Marshal(dto)

			req, _ := http.NewRequest(http.MethodPost, "/driver/locations", bytes.NewBuffer(body))
			req.Header.Add("Content-Type", "application/json")

			jwtToken, _ := createTestToken("valid-token", "Driver-API")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			resp, err := app.Test(req, 20000)
			So(err, ShouldBeNil)

			response := &Respose{}
			json.NewDecoder(resp.Body).Decode(response)
			Data, ok := response.Data.(map[string]interface{})
			So(ok, ShouldBeTrue)

			dtoResp := dtos.DriverBulkCreateResponseDTO{}
			jsonData, _ := json.Marshal(Data)
			json.Unmarshal(jsonData, &dtoResp)

			Convey("Then the handler should be registered and return a match", func() {
				So(response.Status, ShouldEqual, http.StatusCreated)
				So(response.Data, ShouldNotBeNil)
				So(response.Error, ShouldBeEmpty)
				So(response.Meta, ShouldNotBeEmpty)
				So(len(dtoResp.CreatedIDs), ShouldEqual, 2)
			})
		})

		Convey("When registering a handler for GET /driver/locations/nearby", func() {
			app := server.RegisterHandler(http.MethodGet, "/driver/locations/nearby", handler.GetDriversLocationbyNear)
			req, _ := http.NewRequest(http.MethodGet, "/driver/locations/nearby?latitude=40.7128&longitude=-74.0061&radius=10&limit=10", nil)

			jwtToken, _ := createTestToken("valid-token", "Driver-API")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			resp, err := app.Test(req, 20000)
			So(err, ShouldBeNil)

			response := &Respose{}
			json.NewDecoder(resp.Body).Decode(response)
			data, ok := response.Data.(map[string]interface{})
			So(ok, ShouldBeTrue)

			dtoResp := dtos.DriversGetResponseDTO{}
			jsonData, _ := json.Marshal(data)
			json.Unmarshal(jsonData, &dtoResp)

			Convey("Then the handler should be registered and return a match", func() {
				So(response.Status, ShouldEqual, http.StatusOK)
				So(response.Data, ShouldNotBeNil)
				So(response.Error, ShouldBeEmpty)
				So(response.Meta, ShouldNotBeEmpty)
				So(len(dtoResp.Drivers), ShouldBeGreaterThan, 0)
				So(dtoResp.Drivers[0].Distance, ShouldBeGreaterThan, 0)
				So(dtoResp.Drivers[0].ID, ShouldEqual, "1")
				So(dtoResp.Drivers[0].Latitude, ShouldEqual, 40.7128)
				So(dtoResp.Drivers[0].Longitude, ShouldEqual, -74.006)
				So(dtoResp.Drivers[0].Distance, ShouldBeGreaterThan, 0)
			})
		})
	})
}

func createTestToken(tokenString, iss string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"iss": iss,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"sub": "test-user",
	}

	return token.SignedString([]byte(tokenString))
}
