package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bitaksi-case/matching-api/internal/application/dtos"
	"github.com/bitaksi-case/matching-api/internal/application/services"
	"github.com/bitaksi-case/matching-api/internal/errors"
	"github.com/gofiber/fiber/v2"
)

type MatchingHandler struct {
	service *services.MatchingService
	ctx     context.Context
}

func NewMatchingHandler(ctx context.Context, service *services.MatchingService) *MatchingHandler {
	return &MatchingHandler{service: service, ctx: ctx}
}

func (h *MatchingHandler) Match(ctx *fiber.Ctx) error {

	var matchingRequestDTO dtos.MatchingRequestDTO
	if err := ctx.BodyParser(&matchingRequestDTO); err != nil {
		h.sendResponse(ctx, fiber.StatusBadRequest, nil, err)
		return nil
	}

	if err := matchingRequestDTO.Validate(); err != nil {
		h.sendResponse(ctx, fiber.StatusBadRequest, nil, err)
		return nil
	}

	radius, err := strconv.ParseFloat(ctx.Query("radius"), 64)
	if err != nil {
		h.sendResponse(ctx, fiber.StatusBadRequest, nil, err)
		return nil
	}

	limit, err := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if err != nil {
		h.sendResponse(ctx, fiber.StatusBadRequest, nil, err)
		return nil
	}

	matching, err := h.service.Match(h.ctx, &matchingRequestDTO, radius, limit)
	fmt.Println(err)
	if err != nil {
		h.sendResponse(ctx, fiber.StatusInternalServerError, nil, err)
		return nil
	}

	if len(matching.Matching) == 0 {
		h.sendResponse(ctx, fiber.StatusNotFound, nil, errors.ErrLocationNotFound)
		return nil
	}

	h.sendResponse(ctx, fiber.StatusCreated, matching, nil)
	return nil
}

func (h *MatchingHandler) sendResponse(ctx *fiber.Ctx, statusCode int, data interface{}, err error) {
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
	ctx.Set("Content-Type", "application/json")
	ctx.Status(statusCode)
	ctx.Send(jsonResponse)
}
