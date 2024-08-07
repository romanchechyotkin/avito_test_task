//go:build integration

package service

import (
	"context"
	"errors"
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
	defer utest.TeardownTable(log, pg, "flats")

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
		require.ErrorIs(t, err, ErrHouseFlatsNotFound)
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
		require.ErrorIs(t, err, ErrHouseFlatsNotFound)
		require.Equal(t, ([]*entity.Flat)(nil), houseFlats)
	})

}

func TestHouseService_CreateSubscription(t *testing.T) {
	require.NoError(t, prepareErr)

	log.Debug("test configuration", slog.Any("cfg", cfg.Postgresql))

	defer utest.TeardownTable(log, pg, "houses")
	defer utest.TeardownTable(log, pg, "users")

	repositories := repo.NewRepositories(log, pg)

	authService := NewAuthService(log, repositories.User, cfg.JWT.SignKey, cfg.JWT.TokenTTL)
	houseService := NewHouseService(log, repositories.House, repositories.Flat)

	userID, err := authService.CreateUser(context.Background(), &AuthCreateUserInput{
		Email:    "test",
		Password: "123456",
		UserType: "client",
	})
	require.NoError(t, err)

	house, err := houseService.CreateHouse(context.Background(), &HouseCreateInput{
		Address: "Улица Пушкина 1",
		Year:    1999,
	})
	require.NoError(t, err)

	t.Run("successfully created subscription", func(t *testing.T) {
		err = houseService.CreateSubscription(context.Background(), &CreateSubscriptionInput{
			HouseID: fmt.Sprintf("%d", house.ID),
			UserID:  userID,
		})
		require.NoError(t, err)
	})

	t.Run("failed subscription create; house not found", func(t *testing.T) {
		err = houseService.CreateSubscription(context.Background(), &CreateSubscriptionInput{
			HouseID: fmt.Sprintf("%d", house.ID+1),
			UserID:  userID,
		})
		require.ErrorIs(t, err, ErrHouseNotFound)
	})

	t.Run("failed subscription create; subscription already exists", func(t *testing.T) {
		err = houseService.CreateSubscription(context.Background(), &CreateSubscriptionInput{
			HouseID: fmt.Sprintf("%d", house.ID),
			UserID:  userID,
		})
		require.ErrorIs(t, err, ErrHouseSubscriptionExists)
	})

	t.Run("failed subscription create; invalid user id", func(t *testing.T) {
		err = houseService.CreateSubscription(context.Background(), &CreateSubscriptionInput{
			HouseID: fmt.Sprintf("%d", house.ID),
			UserID:  "invalid-user-id", // not uuid
		})
		require.Equal(t, err, err)
		require.False(t, errors.Is(err, ErrHouseNotFound))
		require.False(t, errors.Is(err, ErrHouseSubscriptionExists))
	})
}
