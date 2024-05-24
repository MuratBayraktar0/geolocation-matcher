package dtos

import (
	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/errors"
)

// MatchingRequestDTO swagger:parameters MatchingRequestDTO
type MatchingRequestDTO struct {
	RiderID   string  `json:"rider_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (dto *MatchingRequestDTO) Validate() error {
	if dto.RiderID == "" {
		return errors.ErrInvalidRiderID
	}
	if dto.Latitude < -90 || dto.Latitude > 90 {
		return errors.ErrInvalidLatitude
	}
	if dto.Longitude < -180 || dto.Longitude > 180 {
		return errors.ErrInvalidLongitude
	}

	return nil
}

func (dto *MatchingRequestDTO) ToMatchingEntity() *entities.Location {
	entity := entities.NewLocation(dto.Latitude, dto.Longitude)
	return entity
}

// MatchingResponse swagger:response MatchingResponse
type MatchingResponseDTO struct {
	RiderID  string                 `json:"rider_id"`
	Matching []DriverGetResponseDTO `json:"matching"`
}

type DriverGetResponseDTO struct {
	DriverID  string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
}

func ToDriverGetResponseDTO(matching []*entities.Match) []DriverGetResponseDTO {
	dtos := []DriverGetResponseDTO{}
	for _, match := range matching {
		dto := DriverGetResponseDTO{
			DriverID:  match.DriverID,
			Latitude:  match.Location.Latitude,
			Longitude: match.Location.Longitude,
			Distance:  match.Distance,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}

func ToMatchingResponseDTO(matching *entities.Matching) *MatchingResponseDTO {
	return &MatchingResponseDTO{
		RiderID:  matching.RiderID,
		Matching: ToDriverGetResponseDTO(matching.Matching),
	}
}
