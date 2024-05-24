package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitaksi-case/matching-api/internal/application/dtos"
	"github.com/bitaksi-case/matching-api/internal/application/handlers"
	"github.com/bitaksi-case/matching-api/internal/application/services"
	"github.com/bitaksi-case/matching-api/internal/config"
	"github.com/bitaksi-case/matching-api/internal/domain/usecases"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/repositories"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/server"
	"github.com/dgrijalva/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegration(t *testing.T) {
	Convey("Given a new HTTP server", t, func() {
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
		driverLocationService := services.NewDriverLocationService(httpClientAdapter, authClientAdapter, cfg)
		redisAdapter := adapters.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
		repo := repositories.NewRedisDriverLocationRepository(redisAdapter)
		coreService := usecases.NewMatchingService(driverLocationService, repo)
		service := services.NewMatchingService(coreService)
		handler := handlers.NewMatchingHandler(ctx, service)
		server := server.NewHTTPServer(cfg)
		app := server.RegisterHandler(http.MethodPost, "/matching", handler.Match)

		Convey("When registering a handler for POST /matching", func() {
			dto := dtos.MatchingRequestDTO{
				RiderID:   "rider-1",
				Latitude:  40.712776,
				Longitude: -74.0060,
			}
			body, _ := json.Marshal(dto)

			req, _ := http.NewRequest(http.MethodPost, "/matching?radius=11337&limit=10", bytes.NewBuffer(body))
			req.Header.Add("Content-Type", "application/json")

			jwtToken, _ := createTestToken("valid-token", "User-API")
			req.Header.Set("Authorization", "Bearer "+jwtToken)

			resp, err := app.Test(req, 20000)
			So(err, ShouldBeNil)

			response := &adapters.Respose{}
			json.NewDecoder(resp.Body).Decode(response)
			matchingData, ok := response.Data.(map[string]interface{})
			So(ok, ShouldBeTrue)

			matchingResp := dtos.MatchingResponseDTO{}
			jsonData, _ := json.Marshal(matchingData)
			json.Unmarshal(jsonData, &matchingResp)

			Convey("Then the handler should be registered and return a match", func() {
				So(response.Status, ShouldEqual, http.StatusCreated)
				So(response.Data, ShouldNotBeNil)
				So(response.Error, ShouldBeEmpty)
				So(response.Meta, ShouldNotBeEmpty)
				So(matchingResp.RiderID, ShouldEqual, "rider-1")
				So(len(matchingResp.Matching), ShouldEqual, 2)
				So(matchingResp.Matching[0].DriverID, ShouldEqual, "664bd3cc4384b4cbc5551a50")
				So(matchingResp.Matching[0].Distance, ShouldEqual, 11336386.96747813)
				So(matchingResp.Matching[0].Latitude, ShouldEqual, 41.06040004)
				So(matchingResp.Matching[0].Longitude, ShouldEqual, 27.83085102)
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
