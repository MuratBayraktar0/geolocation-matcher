# Driver Location API

## Overview

This API provides the functionality to find the nearest driver location with a given GeoJSON point and radius. It supports batch operations for bulk updates and includes an endpoint for searching with GeoJSON point and radius fields.

## Getting Started

### Prerequisites

- Go (latest version)
- Docker (if you want to run the application in a container)
- MongoDB

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Run `go mod download` to download the necessary Go modules

### Running the Application

You can run the application in two ways:

#### Directly

Run `go run cmd/main.go`

#### Docker

1. Build the Docker image by running `docker build -t driver-location-api .`
2. Run the Docker container by running `docker run -p 8080:8080 driver-location-api`

## Testing

To run the tests, execute `go test ./...`

## Environment Variables

The application uses the following environment variables, which are defined in `internal/config/local.env` and `internal/config/test.env`:

- `MONGODB_URI`: The URI of your MongoDB instance
- `MONGODB_NAME`: The name of your MongoDB database
- `MONGODB_COLLECTION`: The name of your MongoDB collection
- `SERVER_PORT`: The port on which the server will run
- `SERVER_HOST`: The host of the server
- `SECRET_KEY`: The secret key for token validation
- `ALLOWED_ISSUERS`: The allowed issuers for token validation

## API Endpoints

The API provides the following endpoints:

### POST /driver/locations

Creates multiple driver locations.

Request body:

```json
[
  {
    "id": "string",
    "latitude": "double",
    "longitude": "double"
  }
]
```

### GET /driver/locations/nearby

Returns the drivers locations by near.

Query parameters:

- latitude: double
- longitude: double
- radius: double

```js
GET /driver/locations/nearby?latitude=40.712776&longitude=-74.005974&radius=11337&limit=10
```

For detailed information about each endpoint, refer to the Swagger documentation.

```bash
docker run -p 8080:8080 -e SWAGGER_JSON=/driver-location-api/swagger.json -v /home/user/bitaksi-case/driver-location-api:/driver-location-api swaggerapi/swagger-ui
```

After running the container, you can access the API endpoints at [http://localhost:8080](http://localhost:8080).
