package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
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
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return nil, ErrHouseExists
		}

		s.log.Error("failed to create house in database", logger.Error(err))
		return nil, err
	}

	return house, nil
}

func (s *HouseService) GetHouseFlats(ctx context.Context, input *GetHouseFlatsInput) ([]*entity.Flat, error) {
	flats, err := s.flatRepo.GetHouseFlats(ctx, input.HouseID, input.UserType)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrHouseNotFound
		}

		return nil, err
	}

	return flats, nil
}

func (s *HouseService) CreateSubscription(ctx context.Context, input *CreateSubscriptionInput) error {
	err := s.houseRepo.CreateSubscription(ctx, input.HouseID, input.UserID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return ErrHouseNotFound
		}

		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return ErrHouseSubscriptionExists
		}

		return err
	}

	return nil
}
