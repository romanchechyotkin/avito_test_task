//go:build integration

package service

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/pkg/utest"

	"github.com/stretchr/testify/require"
)

func TestHouseService_CreateHouse(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	defer utest.TeardownTable(log, pg, "houses")

	repositories := repo.NewRepositories(log, pg)

	houseService := NewHouseService(log, repositories.House, repositories.Flat)

	t.Run("successful create house without developer", func(t *testing.T) {
		log.Debug("creating house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address: "Улица Пушкина 1",
			Year:    1999,
		})
		require.NoError(t, err)
		require.Equal(t, "", house.Developer.String)
	})

	t.Run("failed create house via unique constraint", func(t *testing.T) {
		log.Debug("creating non unique house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address: "Улица Пушкина 1",
			Year:    2004,
		})
		require.ErrorIs(t, err, ErrHouseExists)
		require.Equal(t, (*entity.House)(nil), house)
	})

	t.Run("creating house with developer", func(t *testing.T) {
		log.Debug("creating new unique house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address:   "Улица Пушкина 2",
			Year:      2004,
			Developer: "OOO builders",
		})
		require.NoError(t, err)
		require.Equal(t, "OOO builders", house.Developer.String)
	})

	t.Run("creating house without address and year", func(t *testing.T) {
		log.Debug("creating new unique house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Developer: "OOO builders",
		})
		require.ErrorIs(t, err, ErrInvalidInputData)
		require.Equal(t, (*entity.House)(nil), house)
	})
}

func TestHouseService_GetHouseFlats(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	defer utest.TeardownTable(log, pg, "houses")

	repositories := repo.NewRepositories(log, pg)

	flatService := NewFlatService(log, NewSenderService(log, repositories.House), repositories.Flat)
	houseService := NewHouseService(log, repositories.House, repositories.Flat)

	var houseID uint

	t.Run("successful getting flat for house", func(t *testing.T) {
		log.Debug("creating house")
		house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
			Address: "Улица Пушкина 1",
			Year:    1999,
		})
		require.NoError(t, err)
		require.Equal(t, "", house.Developer.String)
		houseID = house.ID

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

		houseFlats, err := houseService.GetHouseFlats(context.Background(), &GetHouseFlatsInput{
			HouseID:  fmt.Sprintf("%d", house.ID),
			UserType: "client",
		})
		require.NoError(t, err)
		require.Equal(t, 0, len(houseFlats))

		houseFlats, err = houseService.GetHouseFlats(context.Background(), &GetHouseFlatsInput{
			HouseID:  fmt.Sprintf("%d", house.ID),
			UserType: "moderator",
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(houseFlats))
	})

	t.Run("failed getting flat for non existing house", func(t *testing.T) {
		houseFlats, err := houseService.GetHouseFlats(context.Background(), &GetHouseFlatsInput{
			HouseID:  fmt.Sprintf("%d", houseID),
			UserType: "client",
		})
		require.ErrorIs(t, err, ErrHouseNotFound)
		require.Equal(t, ([]*entity.Flat)(nil), houseFlats)
	})

}
