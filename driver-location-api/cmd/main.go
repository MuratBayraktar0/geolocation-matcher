package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bitaksi-case/driver-location-api/internal/application/handlers"
	"github.com/bitaksi-case/driver-location-api/internal/application/services"
	"github.com/bitaksi-case/driver-location-api/internal/config"
	"github.com/bitaksi-case/driver-location-api/internal/domain/usecases"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/adapters"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/initdb"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/repositories"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/server"
)

func main() {
	ctx := context.Background()

	// Load configuration
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	cfg := config.LoadConfig(env)

	// Initialize MongoDB client
	mongoAdapter, err := adapters.NewMongoDBClient(ctx, cfg.MongoDBURI)
	if err != nil {
		log.Fatal(err)
	}

	if err = mongoAdapter.Client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Initialize repository
	repo := repositories.NewMongoDBDriverLocationRepository(mongoAdapter, cfg.MongoDBName, cfg.MongoDBCollection)

	// Initialize core service
	coreService := usecases.NewDriverLocationService(repo)

	// Initialize service
	service := services.NewLocationService(coreService)

	// Initialize handlers
	handler := handlers.NewDriverLocationHandler(ctx, service)

	if err := initdb.InitDatabase(ctx, repo); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize HTTP Server
	server := server.NewHTTPServer(cfg)
	server.RegisterHandler(http.MethodPost, "/driver/locations", handler.BulkCreateDriverLocation)
	app := server.RegisterHandler(http.MethodGet, "/driver/locations/nearby", handler.GetDriversLocationbyNear)
	app.Listen(cfg.ServerHost + ":" + cfg.ServerPort)
}
