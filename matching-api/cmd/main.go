package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bitaksi-case/matching-api/internal/application/handlers"
	"github.com/bitaksi-case/matching-api/internal/application/services"
	"github.com/bitaksi-case/matching-api/internal/config"
	"github.com/bitaksi-case/matching-api/internal/domain/usecases"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/repositories"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/server"
)

func main() {
	ctx := context.Background()
	// Load configuration
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	cfg := config.LoadConfig(env)

	redisAdapter := adapters.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	httpClientAdapter := adapters.NewHttpClient(ctx)
	authClientAdapter := adapters.NewAuthClient(cfg.AuthAPIEndpoint)

	// Initialize repository
	repo := repositories.NewRedisDriverLocationRepository(redisAdapter)

	// Initialize driver location service
	driverLocationService := services.NewDriverLocationService(httpClientAdapter, authClientAdapter, cfg)

	// Initialize core service
	coreService := usecases.NewMatchingService(driverLocationService, repo)
	// Initialize service
	service := services.NewMatchingService(coreService)

	// Initialize handlers
	handler := handlers.NewMatchingHandler(ctx, service)

	// Initialize HTTP Server
	server := server.NewHTTPServer(cfg)
	app := server.RegisterHandler(http.MethodPost, "/matching", handler.Match)
	app.Listen(cfg.ServerHost + ":" + cfg.ServerPort)
}
