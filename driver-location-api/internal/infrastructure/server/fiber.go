package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/bitaksi-case/driver-location-api/internal/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/juju/ratelimit"
)

type HTTPServer struct {
	cfg     *config.Config
	app     *fiber.App
	limiter *ratelimit.Bucket
	tokens  int64
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	limiter := ratelimit.NewBucketWithRate(10, 10)
	tokens := limiter.TakeAvailable(3)
	return &HTTPServer{cfg: cfg, app: fiber.New(), limiter: limiter, tokens: tokens}
}

func (s *HTTPServer) RegisterHandler(method, path string, handler fiber.Handler) *fiber.App {
	s.app.Add(method, path, s.RateLimit(s.AuthenticationMiddleware(handler)))
	return s.app
}

func (s *HTTPServer) RateLimit(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		s.tokens = s.limiter.TakeAvailable(3)
		if s.tokens < 3 {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
		}
		return handler(c)
	}
}

func (s *HTTPServer) AuthenticationMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" || !s.isValidToken(token) {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		return next(c)
	}
}

func (s *HTTPServer) isValidToken(token string) bool {
	token = strings.TrimPrefix(token, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		log.Println("Invalid token:", err)
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		log.Println("Invalid token:", err)
		return false
	}

	allowList := s.cfg.AllowedIssuers
	issuer, ok := claims["iss"].(string)
	if !ok || !contains(allowList, issuer) {
		log.Println("Invalid issuer:", issuer)
		return false
	}

	return true
}

func contains(s []string, searchterm string) bool {
	for _, item := range s {
		if item == searchterm {
			return true
		}
	}
	return false
}
