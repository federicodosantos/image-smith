package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/federicodosantos/image-smith/internal/dto"
	model "github.com/federicodosantos/image-smith/internal/model"
	"github.com/federicodosantos/image-smith/internal/usecase"
	customErr "github.com/federicodosantos/image-smith/pkg/error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var CTX = context.TODO()

func createUser() *model.User {
	now := time.Now()
	password := "Rahasia#123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return &model.User{
		ID:        uuid.NewString(),
		Name:      "Jamal",
		Email:     "jamalunyu@gmail.com",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIUserRepository(ctrl)
	mockJWT := NewMockJWTItf(ctrl)

	userUsecase := usecase.NewUserUsecase(mockRepo, mockJWT)

	type testCase struct {
		name             string
		input            *dto.UserRegisterRequest
		mockBehavior     func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf)
		expectedResponse *dto.UserRegisterResponse
		expectError      error
	}

	testCases := []testCase{
		{
			name: "Success - Register new user",
			input: &dto.UserRegisterRequest{
				Name:     "Jamal",
				Email:    "jamalunyu@gmail.com",
				Password: "Rahasia#123",
			},
			mockBehavior: func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf) {
				mockRepo.EXPECT().GetUserByEmail(CTX, "jamalunyu@gmail.com").
					Return(nil, nil)

				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Rahasia#123"), bcrypt.DefaultCost)

				mockRepo.EXPECT().CreateUser(CTX, gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *model.User) error {
						user.ID = uuid.NewString()
						user.CreatedAt = time.Now()
						user.Password = string(hashedPassword)
						return nil
					})
			},
			expectedResponse: &dto.UserRegisterResponse{
				ID:    "mock-uuid",
				Name:  "Jamal",
				Email: "jamalunyu@gmail.com",
			},
			expectError: nil,
		},
		{
			name: "Failed - Password less than 8 letters",
			input: &dto.UserRegisterRequest{
				Name:     "Jamal",
				Email:    "jamalunyu@gmail.com",
				Password: "salah",
			},
			mockBehavior: func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf) {
				mockRepo.EXPECT().GetUserByEmail(CTX, "jamalunyu@gmail.com").
					Return(nil, nil)
			},
			expectedResponse: nil,
			expectError:      fmt.Errorf("Password must be at least 8 characters long"),
		},
		{
			name: "Failed - Email already exists",
			input: &dto.UserRegisterRequest{
				Name:     "Jamal",
				Email:    "jamalunyu@gmail.com",
				Password: "salah",
			},
			mockBehavior: func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf) {
				mockRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "jamalunyu@gmail.com").
					Return(&model.User{}, nil)

				mockRepo.EXPECT().CreateUser(CTX, gomock.Any()).Times(0)
			},
			expectedResponse: nil,
			expectError:      customErr.ErrEmailExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			tc.mockBehavior(mockRepo, mockJWT)

			response, err := userUsecase.Register(ctx, tc.input)

			if tc.expectError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tc.expectedResponse.Name, response.Name)
				assert.Equal(t, tc.expectedResponse.Email, response.Email)
				assert.NotEmpty(t, response.ID)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockIUserRepository(ctrl)
	mockJWT := NewMockJWTItf(ctrl)

	userUsecase := usecase.NewUserUsecase(mockRepo, mockJWT)

	type testCase struct {
		name             string
		input            *dto.UserLoginRequest
		mockBehavior     func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf)
		expectedResponse *dto.UserLoginResponse
		expectError      error
	}

	user := createUser()

	testCases := []testCase{
		{
			name: "Success - Login as a user",
			input: &dto.UserLoginRequest{
				Email:    "jamalunyu@gmail.com",
				Password: "Rahasia#123",
			},
			mockBehavior: func(mockRepo *MockIUserRepository, mockJwt *MockJWTItf) {
				mockRepo.EXPECT().
					GetUserByEmail(CTX, "jamalunyu@gmail.com").
					Return(user, nil)

				mockJWT.EXPECT().CreateToken(user.ID).Return("jwt-token", nil)
			},
			expectedResponse: &dto.UserLoginResponse{
				JWTToken: "jwt-token",
			},
			expectError: nil,
		},
		{
			name: "Failed - Email not found",
			input: &dto.UserLoginRequest{
				Email:    "notfound@gmail.com",
				Password: "Rahasia#123",
			},
			mockBehavior: func(mockRepo *MockIUserRepository, mockJWT *MockJWTItf) {
				mockRepo.EXPECT().
					GetUserByEmail(CTX, "notfound@gmail.com").
					Return(nil, customErr.ErrEmailNotFound)
			},
			expectedResponse: nil,
			expectError:      customErr.ErrEmailNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			tc.mockBehavior(mockRepo, mockJWT)

			response, err := userUsecase.Login(ctx, tc.input)

			if tc.expectError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.JWTToken)
			}
		})
	}
}
