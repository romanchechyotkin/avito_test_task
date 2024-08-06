//go:build integration

package service

import (
	"context"
	"log/slog"
	"testing"

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
		require.Equal(t, nil, house)
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
		require.Error(t, err)
		require.Equal(t, nil, house)
	})
}

func TestHouseService_GetHouseFlats(t *testing.T) {}

func TestHouseService_GetHouse(t *testing.T) {}
