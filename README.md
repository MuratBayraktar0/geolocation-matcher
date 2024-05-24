# Bitaksi Case Study

This repository contains two services: a Driver Location API and a Matching API. The Driver Location API uses location data stored in a MongoDB collection, while the Matching API finds the nearest driver with the rider location using the Driver Location API.

## Getting Started

- Prerequisites
- Docker
- Docker Compose

## Installation

- Clone the repository:
- Navigate to the project directory:
- Build and run the services using Docker Compose:

## Services

### Driver Location API

The Driver Location API provides an endpoint for creating a driver location and an endpoint for searching with GeoJSON point and radius fields. The matched result will have the distance field from the given coordinate.

For more details, refer to the [Driver Location API README](driver-location-api/README.md).

### Matching API

The Matching API matches a suitable driver with the rider using the Driver Location API. It provides an endpoint that allows searching with a GeoJSON point to find a driver if it matches the given criteria.

For more details, refer to the [Driver Location API README](matching-api/README.md).

## Handling High Load and Failures

To handle high load and provide fast response times, we use Redis as a cache. Redis is an in-memory data structure store that can be used as a database, cache, and message broker. It provides a way to quickly access data that would be slower to fetch from MongoDB.

In case of failures or errors, we use a Circuit Breaker mechanism. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail. If the number of failures crosses a certain threshold, the circuit breaker trips, and for the duration of a timeout period, all attempts to invoke the remote service will fail immediately. After the timeout, the circuit breaker allows a limited number of test requests to pass through. If those requests succeed, the circuit breaker resumes normal operation. Otherwise, if there is a failure, the timeout period begins again.

This combination of Redis and Circuit Breaker allows our services to handle high load and prevent failures from cascading and giving the system a chance to recover.

## Geospatial Queries with MongoDB and Redis

Both MongoDB and Redis have built-in support for geospatial data and provide efficient algorithms to perform proximity searches. These algorithms use the Haversine formula to calculate the great-circle distance between two points on a sphere given their longitudes and latitudes, which is the requirement of this case study.

### MongoDB

MongoDB supports geospatial data and queries which allow for queries on data that represents objects defined in a geometric space. For this case, we use the `2dsphere` index and the `$geoNear` operator which returns documents sorted by distance from a specific point.

### Redis

Redis offers data types to handle geospatial data. The `GEOADD` command adds the specified geospatial items (latitude, longitude, name) to the specified key. After adding the geospatial data, we can use the `GEORADIUS` command (referred to as `GeoRadiusQuery` in our code) to query the data, which retrieves items from a specified radius.

By using these built-in geospatial features of MongoDB and Redis, we can efficiently find nearby locations without implementing the Haversine formula manually.
