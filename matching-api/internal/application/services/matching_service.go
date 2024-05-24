package services

import (
	"context"

	"github.com/bitaksi-case/matching-api/internal/application/dtos"
	"github.com/bitaksi-case/matching-api/internal/domain/interfaces"
)

type MatchingService struct {
	service interfaces.MatchingService
}

func NewMatchingService(service interfaces.MatchingService) *MatchingService {
	return &MatchingService{
		service: service,
	}
}

func (s *MatchingService) Match(ctx context.Context, dto *dtos.MatchingRequestDTO, radius float64, limit int64) (*dtos.MatchingResponseDTO, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	location := dto.ToMatchingEntity()
	matching, err := s.service.Matching(ctx, dto.RiderID, location, radius, int64(limit))
	if err != nil {
		return nil, err
	}

	return dtos.ToMatchingResponseDTO(matching), nil
}
