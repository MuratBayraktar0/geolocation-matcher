package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitaksi-case/matching-api/internal/application/dtos"
	"github.com/bitaksi-case/matching-api/internal/config"
	"github.com/bitaksi-case/matching-api/internal/domain/entities"
	"github.com/bitaksi-case/matching-api/internal/domain/interfaces"
	"github.com/bitaksi-case/matching-api/internal/infrastructure/adapters"
	"github.com/sony/gobreaker"
)

type DriverLocationService struct {
	httpClient *adapters.HttpClient
	authClient *adapters.AuthClient
	cfg        *config.Config
	cb         *gobreaker.CircuitBreaker
}

func NewDriverLocationService(httpClient *adapters.HttpClient, authClient *adapters.AuthClient, cfg *config.Config) interfaces.DriverLocationService {
	settings := gobreaker.Settings{
		Name:        "Get drivers location by near",
		MaxRequests: 5,
		Timeout:     10 * time.Second,
	}
	cb := gobreaker.NewCircuitBreaker(settings)

	return &DriverLocationService{
		httpClient: httpClient,
		authClient: authClient,
		cfg:        cfg,
		cb:         cb,
	}
}

func (s *DriverLocationService) GetDriversLocationbyNear(ctx context.Context, point *entities.Location, radius float64, limit int64) (*entities.Drivers, error) {
	token, err := s.authClient.GetToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/driver/locations/nearby?latitude=%f&longitude=%f&radius=%f&limit=%d", s.cfg.DriverLocationApiEndpoint, point.Latitude, point.Longitude, radius, limit)
	resp, err := s.httpClient.Auth(token).Get().CircuitBreaker(s.cb).Request(url)
	if resp != nil && resp.Status == 404 {
		return &entities.Drivers{
			Drivers: []*entities.Driver{},
		}, nil
	}

	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}

	var dto dtos.DriversGetResponseDTO
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		return nil, err
	}

	var driverList []*entities.Driver
	for _, driver := range dto.Drivers {
		driverList = append(driverList, s.toDriverLocationEntity(&driver))
	}

	return &entities.Drivers{
		Drivers: driverList,
	}, nil
}

func (s *DriverLocationService) toDriverLocationEntity(dto *dtos.DriverGetResponseDTO) *entities.Driver {
	return &entities.Driver{
		ID:        dto.DriverID,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
		Distance:  dto.Distance,
	}
}
