package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/internal/service/mocks"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthRoutes_Registration(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *service.AuthCreateUserInput
	}

	type inputBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}

	type MockBehaviour func(m *mocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		inputBody        inputBody
		mockBehavior     MockBehaviour
		wantStatusCode   int
		wantResponseBody string // todo response struct
	}{
		{
			name: "successful registration",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123456",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("test-uuid-id", nil)
			},
			wantStatusCode:   http.StatusCreated,
			wantResponseBody: `{"user_id": "test-uuid-id"}`,
		},
		{
			name: "failed registration; invalid email",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator",
					Password: "123456",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator",
				Password: "123456",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0) // todo return error
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed registration; short password length 3",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed registration; long password length 51",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456789012345678901234567890123456789012345678901",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123456789012345678901234567890123456789012345678901",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "successful registration; min length password length 4",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "1234",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "1234",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("test-uuid-id", nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "successful registration; max length password length 50",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "12345678901234567890123456789012345678901234567890",
					UserType: "moderator",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "12345678901234567890123456789012345678901234567890",
				UserType: "moderator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("test-uuid-id", nil)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "failed registration; email is required",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Password: "123456",
					UserType: "client",
				},
			},
			inputBody: inputBody{
				Password: "123456",
				UserType: "client",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed registration; password is required",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					UserType: "client",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				UserType: "client",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed registration; user type is required",
			args: args{
				ctx: context.Background(),
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mocks.NewMockAuth(ctrl)
			tt.mockBehavior(authService, tt.args)
			services := &service.Services{Auth: authService}

			router := gin.New()
			authGroup := router.Group("/auth")

			newAuthRoutes(logger.NewDiscardLogger(), authGroup, services.Auth)

			reqBody, err := json.Marshal(tt.inputBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", "application-json")

			router.ServeHTTP(recorder, request)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}

func TestAuthRoutes_Login(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *service.AuthGenerateTokenInput
	}

	type inputBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type MockBehaviour func(m *mocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		inputBody        inputBody
		mockBehavior     MockBehaviour
		wantStatusCode   int
		wantResponseBody string // todo response struct
	}{
		{
			name: "successful login",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("test-token", nil)
			},
			wantStatusCode:   http.StatusOK,
			wantResponseBody: `{"token": "test-token"}`,
		},
		{
			name: "failed login; invalid email",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator",
					Password: "123456",
				},
			},
			inputBody: inputBody{
				Email:    "moderator",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0) // todo return error
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; short password length 3",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; long password length 51",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123456789012345678901234567890123456789012345678901",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "123456789012345678901234567890123456789012345678901",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "successful login; min length password length 4",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "1234",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "1234",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("test-token", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "successful login; max length password length 50",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "12345678901234567890123456789012345678901234567890",
				},
			},
			inputBody: inputBody{
				Email:    "moderator@gmail.com",
				Password: "12345678901234567890123456789012345678901234567890",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("test-token", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "failed login; email is required",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Password: "123456",
				},
			},
			inputBody: inputBody{
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; password is required",
			args: args{
				ctx: context.Background(),
				input: &service.AuthGenerateTokenInput{
					Email: "moderator@gmail.com",
				},
			},
			inputBody: inputBody{
				Email: "moderator@gmail.com",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mocks.NewMockAuth(ctrl)
			tt.mockBehavior(authService, tt.args)
			services := &service.Services{Auth: authService}

			router := gin.New()
			authGroup := router.Group("/auth")

			newAuthRoutes(logger.NewDiscardLogger(), authGroup, services.Auth)

			reqBody, err := json.Marshal(tt.inputBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", "application-json")

			router.ServeHTTP(recorder, request)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}
