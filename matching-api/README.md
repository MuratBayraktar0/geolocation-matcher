# Matching API

## Overview

This API provides the functionality to match driver locations with a rider location. It includes an endpoint for matching with a radius and a limit for the number of drivers.

## Getting Started

### Prerequisites

- Go (latest version)
- Docker (if you want to run the application in a container)
- Redis (for caching driver locations)

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Run `go mod download` to download the necessary Go modules

### Running the Application

You can run the application in two ways:

#### Directly

Run `go run cmd/main.go`

#### Docker

1. Build the Docker image by running `docker build -t matching-api .`
2. Run the Docker container by running `docker run -p 8081:8081 matching-api`

## Testing

To run the tests, execute `go test ./...`

## API Endpoints

The API provides the following endpoints:

### POST /matching

Matches driver locations with a rider location.
Request body:

```json
{
    "rider_id": "string",
    "latitude": "double",
    "longitude": "double"
}
```

Query parameters:

- radius: number
- limit: integer

```sh
POST /matching?radius=10&limit=5
```

For detailed information about each endpoint, refer to the Swagger documentation.

```sh
docker run -p 8081:8081 -e SWAGGER_JSON=/matching-api/swagger.json -v /home/user/bitaksi-case/matching-api:/matching-api swaggerapi/swagger-ui
```

After running the container, you can access the API endpoints at [http://localhost:8080](http://localhost:8080).

Please replace /home/user/bitaksi-case/matching-api with the actual path to your project on your machine.
