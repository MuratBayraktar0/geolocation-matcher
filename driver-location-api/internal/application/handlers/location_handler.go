package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitaksi-case/driver-location-api/internal/application/dtos"
	"github.com/bitaksi-case/driver-location-api/internal/application/services"
	"github.com/bitaksi-case/driver-location-api/internal/errors"
	"github.com/gofiber/fiber/v2"
)

type DriverLocationHandler struct {
	service *services.LocationService
	ctx     context.Context
}

func NewDriverLocationHandler(ctx context.Context, service *services.LocationService) *DriverLocationHandler {
	return &DriverLocationHandler{service: service, ctx: ctx}
}

// BulkCreateDriverLocation swagger:route POST /driver/locations drivers bulkCreateDriverLocation
//
// Creates multiple driver locations.
//
// This will create multiple driver locations and return the created driver locations.
//
//	Responses:
//	  201: driverLocationsResponse
func (h *DriverLocationHandler) BulkCreateDriverLocation(c *fiber.Ctx) error {
	dto := new(dtos.DriverBulkCreateRequestDTO)

	if err := c.BodyParser(dto); err != nil {
		h.sendResponse(c, fiber.StatusBadRequest, nil, err)
		return nil
	}
	fmt.Println(dto)
	location, err := h.service.BulkCreateDriverLocation(h.ctx, dto)
	if err != nil {
		if err == errors.ErrInvalidLatitude || err == errors.ErrInvalidLongitude {
			h.sendResponse(c, fiber.StatusBadRequest, nil, err)
			return nil
		}

		h.sendResponse(c, fiber.StatusInternalServerError, nil, err)
		return nil
	}

	h.sendResponse(c, fiber.StatusCreated, location, nil)
	return nil
}

// GetDriversLocationbyNear swagger:route GET /driver/locations/nearby drivers getDriversLocationbyNear
//
// Get drivers locations by near.
//
// This will return the drivers locations by near.
//
//	Responses:
//	  200: driverLocationsResponse
func (h *DriverLocationHandler) GetDriversLocationbyNear(c *fiber.Ctx) error {
	latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		h.sendResponse(c, fiber.StatusBadRequest, nil, errors.ErrInvalidLatitude)
		return nil
	}
	longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		h.sendResponse(c, fiber.StatusBadRequest, nil, errors.ErrInvalidLongitude)
		return nil
	}
	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil {
		h.sendResponse(c, fiber.StatusBadRequest, nil, errors.ErrInvalidRadius)
		return nil
	}
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		h.sendResponse(c, fiber.StatusBadRequest, nil, err)
		return nil
	}

	driversGetRequestDTO := dtos.NewDriverGetRequestDTO(latitude, longitude, radius)
	drivers, err := h.service.GetDriversLocationbyNear(h.ctx, driversGetRequestDTO, limit)
	if err != nil {
		switch err {
		case errors.ErrLocationNotFound:
			h.sendResponse(c, fiber.StatusNotFound, nil, err)
		case errors.ErrInvalidRadius, errors.ErrInvalidLatitude, errors.ErrInvalidLongitude, errors.ErrNilLocation:
			h.sendResponse(c, fiber.StatusBadRequest, nil, err)
		default:
			h.sendResponse(c, fiber.StatusInternalServerError, nil, err)
		}
		return nil
	}

	h.sendResponse(c, fiber.StatusOK, drivers, nil)
	return nil
}

func (h *DriverLocationHandler) sendResponse(c *fiber.Ctx, statusCode int, data interface{}, err error) {
	c.Status(statusCode)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response := map[string]interface{}{
		"status": statusCode,
		"data":   data,
		"error":  nil,
		"meta": map[string]string{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}

	if err != nil {
		response["error"] = err.Error()
		response["data"] = nil
	}

	jsonResponse, _ := json.Marshal(response)
	c.Send(jsonResponse)
}
