package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitaksi-case/driver-location-api/internal/domain/entities"
	"github.com/bitaksi-case/driver-location-api/internal/domain/interfaces"
	"github.com/bitaksi-case/driver-location-api/internal/errors"
	"github.com/bitaksi-case/driver-location-api/internal/infrastructure/adapters"
)

type MongoDBDriverLocationRepository struct {
	client  *mongo.Client
	dbName  string
	colName string
}

func NewMongoDBDriverLocationRepository(client *adapters.MongoDBClient, dbName, colName string) interfaces.DriverLocationRepository {
	return &MongoDBDriverLocationRepository{
		client:  client.Client,
		dbName:  dbName,
		colName: colName,
	}
}

func (repo *MongoDBDriverLocationRepository) BulkCreateLocation(ctx context.Context, drivers *entities.Drivers) (*[]string, error) {
	collection := repo.client.Database(repo.dbName).Collection(repo.colName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var operations []mongo.WriteModel
	for i := 0; i < len(drivers.Drivers); i++ {
		filter := bson.M{"_id": drivers.Drivers[i].ID}
		update := bson.M{
			"$set": drivers.Drivers[i],
		}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	result, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		return nil, err
	}

	var insertedIDs []string
	for _, value := range result.UpsertedIDs {
		insertedIDs = append(insertedIDs, value.(string))
	}
	return &insertedIDs, nil
}

func (repo *MongoDBDriverLocationRepository) GetDriversLocationbyNear(ctx context.Context, loc *entities.Location, radius float64, limit int64) (*entities.Drivers, error) {
	collection := repo.client.Database(repo.dbName).Collection(repo.colName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$geoNear": bson.M{
				"near":          loc,
				"distanceField": "distance",
				"maxDistance":   radius * 1000,
				"spherical":     true,
			},
		},
		{
			"$limit": limit,
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []*entities.Driver
	for cursor.Next(ctx) {
		var doc entities.Driver
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		drivers = append(drivers, &doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		return nil, errors.ErrLocationNotFound
	}

	return entities.NewDrivers(drivers), nil
}

func (repo *MongoDBDriverLocationRepository) GetLocationCount(ctx context.Context) (int64, error) {
	collection := repo.client.Database(repo.dbName).Collection(repo.colName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *MongoDBDriverLocationRepository) CreateLocationIndex(ctx context.Context) error {
	collection := repo.client.Database(repo.dbName).Collection(repo.colName)
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"location": "2dsphere", // Create a 2dsphere index on the "location" field
		},
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
