package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPServer_RegisterHandler(t *testing.T) {
	Convey("Given a new HTTP server", t, func() {
		cfg := config.LoadConfig("test")
		server := NewHTTPServer(cfg)

		Convey("When registering a handler for GET /hello", func() {
			handler := func(ctx *fiber.Ctx) error {
				ctx.SendStatus(fiber.StatusOK)
				return nil
			}
			app := server.RegisterHandler("GET", "/hello", handler)

			Convey("Then the handler should be registered)", func() {
				req, _ := http.NewRequest("GET", "/hello", nil)
				token, _ := createTestToken(cfg.SecretKey, "Matching-API")
				req.Header.Set("Authorization", "Bearer "+token)

				resp, _ := app.Test(req)
				So(resp.StatusCode, ShouldEqual, fiber.StatusOK)
			})
		})
	})
}

func createTestToken(secretKey string, issuer string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":     issuer,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		issuer:    issuer,
		secretKey: secretKey,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
