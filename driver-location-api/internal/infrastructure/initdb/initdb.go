package initdb

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InitDatabase initializes the database with initial data from a CSV file.
func InitDatabase(ctx context.Context, repo interfaces.DriverLocationRepository) error {
	count, _ := repo.GetLocationCount(ctx)
	if count > 0 {
		return nil
	}

	file, err := os.Open("data/Coordinates.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	chunkSize := 10000 // Set the desired chunk size
	var drivers []*entities.Driver

	errChan := make(chan error)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		latitude, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			continue
		}

		longitude, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		driver := entities.NewDriver(primitive.NewObjectID().Hex(), entities.NewLocation(latitude, longitude), time.Now(), time.Now())
		drivers = append(drivers, driver)

		if len(drivers) >= chunkSize {
			go func(chunk []*entities.Driver) {
				if _, err := repo.BulkCreateLocation(ctx, entities.NewDrivers(chunk)); err != nil {
					errChan <- err
				}
			}(drivers)
			drivers = nil // reset the drivers slice
		}
	}

	// handle any remaining drivers
	if len(drivers) > 0 {
		go func(chunk []*entities.Driver) {
			if _, err := repo.BulkCreateLocation(ctx, entities.NewDrivers(chunk)); err != nil {
				errChan <- err
			}
		}(drivers)
	}

	select {
	case err := <-errChan:
		return err
	default:
		repo.CreateLocationIndex(ctx)
		return nil
	}
}
