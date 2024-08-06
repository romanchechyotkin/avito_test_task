//go:build integration

package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/utest"

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
