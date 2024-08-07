package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"log/slog"
)

type FlatService struct {
	log *slog.Logger

	sendService Sender
	flatRepo    repo.Flat
}

func NewFlatService(log *slog.Logger, sendService Sender, flatRepo repo.Flat) *FlatService {
	log = log.With(slog.String("component", "flat service"))

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
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrHouseNotFound
		}

		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return nil, ErrFlatExists
		}

		s.log.Debug("failed to create flat in database", logger.Error(err))
		return nil, err
	}

	s.log.Info("created new flat", slog.Any("flat", flat))

	return flat, nil
}

func (s *FlatService) UpdateFlat(ctx context.Context, input *FlatUpdateInput) (*entity.Flat, error) {
	status, err := s.flatRepo.GetStatus(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrFlatNotFound
		}

		return nil, err
	}

	if status == "created" && input.Status != "on moderation" {
		return nil, ErrFlatNotOnModeration
	}

	if status == "on moderation" && input.Status == "on moderation" {
		return nil, ErrFlatOnModeration
	}

	if status == "approved" || status == "declined" {
		return nil, ErrFlatFinishedModeration
	}

	if input.Status == "created" || input.Status == "approved" || input.Status == "declined" {
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

	if flat.ModerationStatus == "approved" {
		s.log.Info("sending email to subscribers", slog.Any("house id", flat.HouseID))
		go func() {
			s.sendService.Send() <- flat.HouseID
		}()
	}

	s.log.Info("updated flat", slog.Any("id", flat.ID), slog.String("status", flat.ModerationStatus))

	return flat, nil
}
