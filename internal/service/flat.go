package service

import (
	"context"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
)

type FlatService struct {
	log *slog.Logger

	flatRepo repo.Flat
}

func NewFlatService(log *slog.Logger, flatRepo repo.Flat) *FlatService {
	return &FlatService{
		log:      log,
		flatRepo: flatRepo,
	}
}

func (s *FlatService) CreateFlat(ctx context.Context, input *FlatCreateInput) (*entity.Flat, error) {
	flat, err := s.flatRepo.CreateFlat(ctx, &entity.Flat{
		Number:      input.Number,
		HouseID:     input.HouseID,
		Price:       input.Price,
		RoomsAmount: input.RoomsAmount,
	})
	if err != nil {
		return nil, err
	}

	return flat, nil
}
