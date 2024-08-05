package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

func TestHouseRoutes_CreateHouse(t *testing.T) {
	type args struct {
		userType         string
		userID           string
		createHouseInput *service.HouseCreateInput
		houseEntity      *entity.House
		isAuth           bool
		token            string
	}

	type AuthMockBehaviour func(m *mocks.MockAuth, args args)
	type HouseMockBehaviour func(m *mocks.MockHouse, args args)

	testCases := []struct {
		name              string
		args              args
		reqBody           request.CreateHouse
		houseMockBehavior HouseMockBehaviour
		authMockBehavior  AuthMockBehaviour
		wantStatusCode    int
	}{
		{
			name: "successful create",
			reqBody: request.CreateHouse{
				Address:   "Ул Пушкина 1",
				Year:      1999,
				Developer: "",
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
					Year:    1999,
				},
				houseEntity: &entity.House{
					Address: "Ул Пушкина 1",
					Year:    1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "failed create; client user",
			reqBody: request.CreateHouse{
				Address:   "Ул Пушкина 1",
				Year:      1999,
				Developer: "",
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
					Year:    1999,
				},
				houseEntity: &entity.House{
					Address: "Ул Пушкина 1",
					Year:    1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "failed create; no authorization",
			reqBody: request.CreateHouse{
				Address:   "Ул Пушкина 1",
				Year:      1999,
				Developer: "",
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
					Year:    1999,
				},
				houseEntity: &entity.House{
					Address: "Ул Пушкина 1",
					Year:    1999,
				},
				userType: "client",
				userID:   "test-uuid-id",
				isAuth:   false,
			},
			authMockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().ParseToken(gomock.Any()).Return(&service.TokenClaims{
					UserType: args.userType,
					UserID:   args.userID,
				}, nil).Times(0)
			},
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed create; invalid authorization",
			reqBody: request.CreateHouse{
				Address:   "Ул Пушкина 1",
				Year:      1999,
				Developer: "",
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
					Year:    1999,
				},
				houseEntity: &entity.House{
					Address: "Ул Пушкина 1",
					Year:    1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed create; address is required",
			reqBody: request.CreateHouse{
				Year: 1999,
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Year: 1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; year is required",
			reqBody: request.CreateHouse{
				Address: "Ул Пушкина 1",
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(args.houseEntity, nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; house exists",
			reqBody: request.CreateHouse{
				Address: "Ул Пушкина 1",
				Year:    1999,
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 1",
					Year:    1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(nil, service.ErrHouseExists)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed create; internal server error",
			reqBody: request.CreateHouse{
				Address: "Ул Пушкина 123",
				Year:    1999,
			},
			args: args{
				createHouseInput: &service.HouseCreateInput{
					Address: "Ул Пушкина 123",
					Year:    1999,
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
			houseMockBehavior: func(m *mocks.MockHouse, args args) {
				m.EXPECT().CreateHouse(gomock.Any(), args.createHouseInput).Return(nil, errors.New("some error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			houseService := mocks.NewMockHouse(ctrl)
			tt.houseMockBehavior(houseService, tt.args)

			authService := mocks.NewMockAuth(ctrl)
			tt.authMockBehavior(authService, tt.args)

			services := &service.Services{House: houseService}

			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = io.Discard
			router := gin.New()
			houseGroup := router.Group("/v1/house")

			authMiddleware := middleware.NewAuthMiddleware(authService)

			newHouseRoutes(logger.NewDiscardLogger(), houseGroup, services.House, authMiddleware)

			reqBody, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/house/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application-json")

			if tt.args.isAuth {
				req.Header.Set("Authorization", tt.args.token)
			}

			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}

func TestHouseRoutes_GetHouseFlats(t *testing.T) {

}
