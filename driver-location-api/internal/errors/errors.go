package errors

import "errors"

var (
	// ErrLocationNotFound is returned when a location is not found in the database.
	ErrLocationNotFound = errors.New("location not found")

	// ErrInvalidLocation is returned when the provided location data is invalid.
	ErrInvalidLocation = errors.New("invalid location data")

	// ErrNilLocation is returned when the provided location is nil.
	ErrNilLocation = errors.New("location cannot be nil")

	// ErrInvalidRadius is returned when the provided radius is invalid.
	ErrInvalidRadius = errors.New("radius must be greater than zero")

	// ErrInvalidLatitude is returned when the provided latitude is invalid.
	ErrInvalidLatitude = errors.New("latitude must be between -90 and 90")

	// ErrInvalidLongitude is returned when the provided longitude is invalid.
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")

	// ErrInvalidDriverID is returned when the provided driver ID is invalid.
	ErrInvalidDriverID = errors.New("driver ID must be greater than zero")

	// ErrObjectIDConversion is returned when the InsertedID cannot be converted to ObjectID.
	ErrObjectIDConversion = errors.New("failed to convert InsertedID to ObjectID")
)
