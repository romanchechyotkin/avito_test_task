package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
)

type HouseService struct {
	log *slog.Logger

	houseRepo repo.House
}

func NewHouseService(log *slog.Logger, houseRepo repo.House) *HouseService {
	return &HouseService{
		log:       log,
		houseRepo: houseRepo,
	}
}

func (s *HouseService) CreateHouse(ctx context.Context, input *HouseCreateInput) (*entity.House, error) {
	house, err := s.houseRepo.CreateHouse(ctx, &entity.House{
		Address: input.Address,
		Year:    input.Year,
		Developer: sql.NullString{
			String: input.Developer,
		},
	})
	if err != nil {
		return nil, err
	}

	return house, nil
}
