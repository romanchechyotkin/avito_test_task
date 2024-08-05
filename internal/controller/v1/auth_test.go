package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/internal/service/mocks"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthRoutes_Registration(t *testing.T) {
	type args struct {
		input *service.AuthCreateUserInput
	}

	type MockBehaviour func(m *mocks.MockAuth, args args)

	testCases := []struct {
		name             string
		args             args
		reqBody          request.Registration
		mockBehavior     MockBehaviour
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "successful registration",
			args: args{
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator",
					Password: "123456",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456789012345678901234567890123456789012345678901",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "1234",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "12345678901234567890123456789012345678901234567890",
					UserType: "moderator",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Password: "123456",
					UserType: "client",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					UserType: "client",
				},
			},
			reqBody: request.Registration{
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
				input: &service.AuthCreateUserInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
				},
			},
			reqBody: request.Registration{
				Email:    "moderator@gmail.com",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed registration; invalid user type",
			args: args{
				input: &service.AuthCreateUserInput{},
			},
			reqBody: request.Registration{
				Email:    "operator@gmail.com",
				Password: "123456",
				UserType: "operator",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().CreateUser(gomock.Any(), args.input).Return("", nil).Times(0) // todo return error
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

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application-json")

			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}

func TestAuthRoutes_Login(t *testing.T) {
	type args struct {
		input *service.AuthGenerateTokenInput
	}

	type MockBehaviour func(m *mocks.MockAuth, args args)

	testCases := []struct {
		name           string
		args           args
		reqBody        request.Login
		mockBehavior   MockBehaviour
		wantStatusCode int
	}{
		{
			name: "successful login",
			args: args{
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123456",
				},
			},
			reqBody: request.Login{
				Email:    "moderator@gmail.com",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("test-token", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "failed login; invalid email",
			args: args{
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator",
					Password: "123456",
				},
			},
			reqBody: request.Login{
				Email:    "moderator",
				Password: "123456",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; short password length 3",
			args: args{
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123",
				},
			},
			reqBody: request.Login{
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
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123456789012345678901234567890123456789012345678901",
				},
			},
			reqBody: request.Login{
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
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "1234",
				},
			},
			reqBody: request.Login{
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
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "12345678901234567890123456789012345678901234567890",
				},
			},
			reqBody: request.Login{
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
				input: &service.AuthGenerateTokenInput{
					Password: "123456",
				},
			},
			reqBody: request.Login{
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
				input: &service.AuthGenerateTokenInput{
					Email: "moderator@gmail.com",
				},
			},
			reqBody: request.Login{
				Email: "moderator@gmail.com",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", nil).Times(0)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; wrong password",
			args: args{
				input: &service.AuthGenerateTokenInput{
					Email:    "moderator@gmail.com",
					Password: "123123213",
				},
			},
			reqBody: request.Login{
				Email:    "moderator@gmail.com",
				Password: "123123213",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", service.ErrWrongPassword)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed login; wrong email",
			args: args{
				input: &service.AuthGenerateTokenInput{
					Email:    "123@gmail.com",
					Password: "12312321313",
				},
			},
			reqBody: request.Login{
				Email:    "123@gmail.com",
				Password: "12312321313",
			},
			mockBehavior: func(m *mocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(gomock.Any(), args.input).Return("", service.ErrUserNotFound)
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

			body, err := json.Marshal(tt.reqBody)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application-json")

			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantStatusCode, recorder.Code)
		})
	}
}
