package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/middleware"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/entity"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/internal/service/mocks"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFlatRoutes_CreateFlat(t *testing.T) {
	type args struct {
		userType        string
		userID          string
		createFlatInput *service.FlatCreateInput
		flatEntity      *entity.Flat
		isAuth          bool
		token           string
	}

	type AuthMockBehaviour func(m *mocks.MockAuth, args args)
	type FlatMockBehaviour func(m *mocks.MockFlat, args args)

	testCases := []struct {
		name              string
		args              args
		reqBody           request.CreateFlat
		authMockBehavior  AuthMockBehaviour
		flatMockBehaviour FlatMockBehaviour
		wantStatusCode    int
	}{
		{
			name: "successful create by moderator",
			reqBody: request.CreateFlat{
				Number:      1,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      1,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           1,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "moderator",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "successful create by client",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "failed create; no authorization",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   false,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil).Times(0)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed create; invalid token",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil).Times(0)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed create; no user type",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil).Times(0)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed create; number is required",
			reqBody: request.CreateFlat{
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{},
				flatEntity:      &entity.Flat{},
				userType:        "moderator",
				userID:          "test-uuid-id",
				isAuth:          true,
				token:           "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; house id is required",
			reqBody: request.CreateFlat{
				Number:      1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{},
				flatEntity:      &entity.Flat{},
				userType:        "moderator",
				userID:          "test-uuid-id",
				isAuth:          true,
				token:           "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; price is required",
			reqBody: request.CreateFlat{
				Number:      1,
				HouseID:     1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{},
				flatEntity:      &entity.Flat{},
				userType:        "moderator",
				userID:          "test-uuid-id",
				isAuth:          true,
				token:           "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; rooms amount is required",
			reqBody: request.CreateFlat{
				Number:  1,
				HouseID: 1,
				Price:   1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{},
				flatEntity:      &entity.Flat{},
				userType:        "moderator",
				userID:          "test-uuid-id",
				isAuth:          true,
				token:           "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(args.flatEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; house not found",
			reqBody: request.CreateFlat{
				Number:      1,
				HouseID:     11,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      1,
					HouseID:     11,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           1,
					HouseID:          11,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(nil, service.ErrHouseNotFound)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; flat exists",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(nil, service.ErrFlatExists)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; internal server error",
			reqBody: request.CreateFlat{
				Number:      2,
				HouseID:     1,
				Price:       1,
				RoomsAmount: 1,
			},
			args: args{
				createFlatInput: &service.FlatCreateInput{
					Number:      2,
					HouseID:     1,
					Price:       1,
					RoomsAmount: 1,
				},
				flatEntity: &entity.Flat{
					Number:           2,
					HouseID:          1,
					Price:            1,
					RoomsAmount:      1,
					ModerationStatus: "created",
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   true,
				token:    "Bearer test-token",
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil)
			},
			flatMockBehaviour: func(m *mocks.MockFlat, args args) {
				m.EXPECT().CreateFlat(gomock.Any(), args.createFlatInput).Return(nil, errors.New("some error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			flatService := mocks.NewMockFlat(ctrl)
			tt.flatMockBehaviour(flatService, tt.args)

			authService := mocks.NewMockAuth(ctrl)
			tt.authMockBehavior(authService, tt.args)

			services := &service.Services{Flat: flatService}

			router := gin.New()
			houseGroup := router.Group("/v1/flat")

			authMiddleware := middleware.NewAuthMiddleware(authService)

			newFlatRoutes(logger.NewDiscardLogger(), houseGroup, services.Flat, authMiddleware)

			reqBody, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/flat/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application-json")

			if tt.args.isAuth {
				req.Header.Set("Authorization", tt.args.token)
			}

			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}
