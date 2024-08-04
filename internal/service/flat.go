package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
)

type FlatService struct {
	log *slog.Logger

	sendService Sender
	flatRepo    repo.Flat
}

func NewFlatService(log *slog.Logger, sendService Sender, flatRepo repo.Flat) *FlatService {
	return &FlatService{
		log:         log,
		sendService: sendService,
		flatRepo:    flatRepo,
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

func (s *FlatService) UpdateFlat(ctx context.Context, input *FlatUpdateInput) (*entity.Flat, error) {
	status, err := s.flatRepo.GetStatus(ctx, input.ID)
	if err != nil {
		// todo NoRowsError
		return nil, err
	}

	// todo custom errors

	if status == "created" && input.Status != "on moderation" {
		return nil, errors.New("сначала надо взять квартиру на модерацию")
	}

	if status == "on moderation" && input.Status == "on moderation" {
		return nil, errors.New("квартира уже на модерарации")
	}

	if status == "approved" || status == "declined" {
		return nil, errors.New("квартира уже прошла модерарацию")
	}

	if input.Status == "created" {
		input.ModeratorID = ""
	}

	if input.Status == "approved" || input.Status == "declined" {
		input.ModeratorID = ""
	}

	flat, err := s.flatRepo.UpdateStatus(ctx, &entity.Flat{
		ID:               input.ID,
		ModerationStatus: input.Status,
	}, sql.NullString{
		String: input.ModeratorID,
		Valid:  len(input.ModeratorID) > 0,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		s.sendService.Notify() <- flat.HouseID
	}()

	return flat, nil
}
