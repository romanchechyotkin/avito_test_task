package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/repo"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID   uint
	UserType string
}

type AuthService struct {
	log *slog.Logger

	signKey  string
	tokenTTL time.Duration
	userRepo repo.User
}

func NewAuthService(log *slog.Logger, userRepo repo.User, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:      log,
		signKey:  signKey,
		tokenTTL: tokenTTL,
		userRepo: userRepo,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input *AuthCreateUserInput) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}

	user := &entity.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		UserType: input.UserType,
	}

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrors.ErrUserExists) {
			return 0, err
		}

		s.log.Error("failed to create user in database", logger.Error(err))
		return 0, err
	}

	return userID, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input *AuthGenerateTokenInput) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return "", err
		}

		s.log.Error("failed to get user by email", logger.Error(err))
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", ErrWrongPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID:   user.ID,
		UserType: user.UserType,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		s.log.Error("failed to sign token", logger.Error(err))

		return "", ErrSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accesstoken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return nil, ErrParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrParseToken
	}

	return claims, nil
}
