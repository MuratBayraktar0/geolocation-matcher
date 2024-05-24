package adapters_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/dgrijalva/jwt-go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/sony/gobreaker"
)

type MockRoundTripper struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestHttpClient(t *testing.T) {
	Convey("Given a new HttpClient", t, func() {
		ctx := context.Background()
		client := adapters.NewHttpClient(ctx)

		Convey("When setting the authorization token", func() {
			token := "my-token"
			client.Auth(token)

			Convey("The headers should contain the authorization token", func() {
				So(client.Headers["Authorization"], ShouldResemble, []string{"Bearer " + token})
			})
		})

		Convey("When setting the HTTP method to POST", func() {
			client.Post()

			Convey("The method should be set to POST", func() {
				So(client.Method, ShouldEqual, http.MethodPost)
			})
		})

		Convey("When setting the HTTP method to GET", func() {
			client.Get()

			Convey("The method should be set to GET", func() {
				So(client.Method, ShouldEqual, http.MethodGet)
			})
		})

		Convey("When setting the CircuitBreaker", func() {
			cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{})

			client.CircuitBreaker(cb)

			Convey("The CircuitBreaker should be set", func() {
				So(client.CB, ShouldEqual, cb)
			})
		})
	})
}

func TestHttpClient_Request(t *testing.T) {
	Convey("Given a new HttpClient", t, func() {
		ctx := context.Background()
		client := adapters.NewHttpClient(ctx)
		jwtToken, _ := createTestToken("valid-token", "User-API")

		Convey("When making a GET request", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := map[string]interface{}{
					"status": 200,
					"data":   "example",
				}

				jsonResponse, _ := json.Marshal(response)
				w.Write(jsonResponse)
			}))
			defer ts.Close()
			resp, err := client.Auth(jwtToken).Get().Request(ts.URL)
			Convey("The request should be sent successfully", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.Status, ShouldEqual, 200)
				So(resp.Data, ShouldEqual, "example")
			})
		})

		Convey("When making a POST request", func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := map[string]interface{}{
					"status": 201,
					"data":   "example",
				}

				jsonResponse, _ := json.Marshal(response)
				w.Write(jsonResponse)
			}))
			defer ts.Close()
			resp, err := client.Auth(jwtToken).Post().Request(ts.URL)
			Convey("The request should be sent successfully", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.Status, ShouldEqual, 201)
				So(resp.Data, ShouldEqual, "example")
			})
		})
	})
}

func TestHttpClient_Request_CircuitBreaker(t *testing.T) {
	Convey("Given a new HttpClient with a CircuitBreaker", t, func() {
		ctx := context.Background()
		client := adapters.NewHttpClient(ctx)
		cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{})

		client.CircuitBreaker(cb)

		Convey("When making a request that fails", func() {
			url := "https://example.com"
			client.Get()

			Convey("The request should be sent successfully", func() {
				// Set the mock client
				mockClient := &MockRoundTripper{
					DoFunc: func(req *http.Request) (*http.Response, error) {
						So(req.Method, ShouldEqual, http.MethodPost)
						So(req.URL.String(), ShouldEqual, url)
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Body:       io.NopCloser(strings.NewReader(`{"status": 404, "error": "Not Found"}`)),
						}, nil
					},
				}

				client.Client = &http.Client{Transport: mockClient}
				for i := 0; i < 15; i++ {
					client.Get()

					// Make the request
					resp, err := client.Request(url)

					// Verify the response
					So(err, ShouldNotBeNil)
					So(resp, ShouldBeNil)
				}
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
