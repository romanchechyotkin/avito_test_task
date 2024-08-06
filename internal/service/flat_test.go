//go:build integration

package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/utest"

	"github.com/docker/distribution/uuid"
	"github.com/stretchr/testify/require"
)

func TestFlatService_CreateFlat(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	defer utest.TeardownTable(log, pg, "houses")

	repositories := repo.NewRepositories(log, pg)

	houseService := NewHouseService(log, repositories.House, repositories.Flat)
	flatService := NewFlatService(log, NewSenderService(log, repositories.House), repositories.Flat)

	var houseID uint

	t.Run("successful create flat for house", func(t *testing.T) {
		log.Debug("creating house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address: "Улица Пушкина 1",
			Year:    1999,
		})
		require.NoError(t, err)
		require.Equal(t, "", house.Developer.String)

		log.Debug("creating flat")
		flat, err := flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      1,
			HouseID:     house.ID,
			Price:       1,
			RoomsAmount: 1,
		})
		require.NoError(t, err)
		require.Equal(t, house.ID, flat.HouseID)
		require.Equal(t, "created", flat.ModerationStatus)
		houseID = flat.HouseID
	})

	t.Run("failed create flat for non existing house", func(t *testing.T) {
		log.Debug("creating flat")
		flat, err := flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      1,
			HouseID:     2,
			Price:       1,
			RoomsAmount: 1,
		})
		require.ErrorIs(t, err, ErrHouseNotFound)
		require.Equal(t, (*entity.Flat)(nil), flat)
	})

	t.Run("failed create flat via unique constraint", func(t *testing.T) {
		log.Debug("creating flat")
		flat, err := flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      1,
			HouseID:     houseID,
			Price:       1,
			RoomsAmount: 1,
		})
		require.ErrorIs(t, err, ErrFlatExists)
		require.Equal(t, (*entity.Flat)(nil), flat)
	})

	t.Run("successful create flats with same number for different houses", func(t *testing.T) {
		log.Debug("creating house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address:   "Улица Пушкина 2",
			Year:      2005,
			Developer: "OOO builders",
		})
		require.NoError(t, err)
		require.Equal(t, "OOO builders", house.Developer.String)

		log.Debug("creating flat #2")
		flat, err := flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      1,
			HouseID:     house.ID,
			Price:       1,
			RoomsAmount: 1,
		})
		require.NoError(t, err)
		require.Equal(t, house.ID, flat.HouseID)

		log.Debug("creating flat #3")
		flat, err = flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      2,
			HouseID:     houseID,
			Price:       1,
			RoomsAmount: 1,
		})
		require.NoError(t, err)
		require.Equal(t, houseID, flat.HouseID)

		log.Debug("creating flat #4")
		flat, err = flatService.CreateFlat(context.Background(), &FlatCreateInput{
			Number:      2,
			HouseID:     house.ID,
			Price:       1,
			RoomsAmount: 1,
		})
		require.NoError(t, err)
		require.Equal(t, house.ID, flat.HouseID)
	})
}

func TestFlatService_UpdateFlat(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	defer utest.TeardownTable(log, pg, "houses")

	repositories := repo.NewRepositories(log, pg)

	houseService := NewHouseService(log, repositories.House, repositories.Flat)
	flatService := NewFlatService(log, NewSenderService(log, repositories.House), repositories.Flat)

	modearorID := uuid.Generate().String()

	house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
		Address: "Улица Пушкина 1",
		Year:    1999,
	})
	require.NoError(t, err)
	require.Equal(t, "", house.Developer.String)

	log.Debug("creating flat")
	flat, err := flatService.CreateFlat(context.Background(), &FlatCreateInput{
		Number:      1,
		HouseID:     house.ID,
		Price:       1,
		RoomsAmount: 1,
	})
	require.NoError(t, err)
	require.Equal(t, house.ID, flat.HouseID)
	require.Equal(t, "created", flat.ModerationStatus)

	t.Run("failed update flat for non existing flat", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID + 1,
			Status:      "on moderation",
			ModeratorID: modearorID,
		})
		require.ErrorIs(t, err, ErrFlatNotFound)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

	t.Run("failed update flat; wrong moderator id type", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "on moderation",
			ModeratorID: "id", // not uuid
		})
		require.Error(t, err)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

	t.Run("failed update flat; approved status not allowed from created", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "approved",
			ModeratorID: modearorID,
		})
		require.Error(t, err, ErrFlatNotOnModeration)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

	t.Run("failed update flat; declined status not allowed from created", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "declined",
			ModeratorID: modearorID,
		})
		require.Error(t, err, ErrFlatNotOnModeration)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

	t.Run("successful update; from created to on moderation", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "on moderation",
			ModeratorID: modearorID,
		})
		require.NoError(t, err)
		require.Equal(t, "on moderation", updatedFlat.ModerationStatus)
	})

	t.Run("failed update flat; on moderation status not allowed from on moderation", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "on moderation",
			ModeratorID: modearorID,
		})
		require.Error(t, err, ErrFlatOnModeration)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

	t.Run("successful update flat; from on moderation status to created", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "created",
			ModeratorID: modearorID,
		})
		require.NoError(t, err)
		require.Equal(t, "created", updatedFlat.ModerationStatus)
	})

	t.Run("successful update; from created to on moderation to approved", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "on moderation",
			ModeratorID: modearorID,
		})
		require.NoError(t, err)
		require.Equal(t, "on moderation", updatedFlat.ModerationStatus)

		updatedFlat, err = flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "approved",
			ModeratorID: modearorID,
		})
		require.NoError(t, err)
		require.Equal(t, "approved", updatedFlat.ModerationStatus)
	})

	t.Run("failed update flat; not allowed moderate flat when its already approved or declined", func(t *testing.T) {
		updatedFlat, err := flatService.UpdateFlat(context.Background(), &FlatUpdateInput{
			ID:          flat.ID,
			Status:      "declined",
			ModeratorID: modearorID,
		})
		require.ErrorIs(t, err, ErrFlatFinishedModeration)
		require.Equal(t, (*entity.Flat)(nil), updatedFlat)
	})

}
