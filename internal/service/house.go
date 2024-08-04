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
	flatRepo  repo.Flat
}

func NewHouseService(log *slog.Logger, houseRepo repo.House, flatRepo repo.Flat) *HouseService {
	return &HouseService{
		log:       log,
		houseRepo: houseRepo,
		flatRepo:  flatRepo,
	}
}

func (s *HouseService) CreateHouse(ctx context.Context, input *HouseCreateInput) (*entity.House, error) {
	house, err := s.houseRepo.CreateHouse(ctx, &entity.House{
		Address: input.Address,
		Year:    input.Year,
		Developer: sql.NullString{
			Valid:  len(input.Developer) > 0,
			String: input.Developer,
		},
	})
	if err != nil {
		return nil, err
	}

	return house, nil
}

func (s *HouseService) GetHouseFlats(ctx context.Context, input *GetHouseFlatsInput) ([]*entity.Flat, error) {
	flats, err := s.flatRepo.GetHouseFlats(ctx, input.HouseID, input.UserType)
	if err != nil {
		return nil, err
	}

	return flats, nil
}
