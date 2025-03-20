package delivery_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/federicodosantos/image-smith/internal/delivery"
	"github.com/federicodosantos/image-smith/internal/dto"
	customErr "github.com/federicodosantos/image-smith/pkg/error"
	response "github.com/federicodosantos/image-smith/pkg/response"
	"go.uber.org/mock/gomock"
)

var jsonHeader = http.Header{
	"Content-Type": {"application/json"}}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockIUserUsecase(ctrl)
	userHandler := delivery.NewUserHandler(mockUsecase)

	type parameter struct {
		w http.ResponseWriter
		r *http.Request
	}

	type TestCase struct {
		Name           string
		Input          parameter
		mockBehavior   func(mockUsecase *MockIUserUsecase)
		expectedHeader http.Header
		expectedBody   response.HttpResponse
	}

	testCases := []TestCase{
		{
			Name: "Success - Register a new user",
			Input: parameter{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://0.0.0.0/auth/register",
					strings.NewReader(`
						{
							"name":     "Jamal",
							"email":    "jamalunyu@gmail.com",
							"password": "Rahasia#123"
						},
					`),
				),
			},
			mockBehavior: func(mockUsecase *MockIUserUsecase) {
				mockUsecase.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(&dto.UserRegisterResponse{
						ID:        "uuid",
						Name:      "Jamal",
						Email:     "jamalunyu@gmail.com",
						CreatedAt: time.Now(),
					}, nil)
			},
			expectedHeader: jsonHeader,
			expectedBody: response.HttpResponse{
				Status:  http.StatusCreated,
				Message: "successfully create user",
				Data: &dto.UserRegisterResponse{
					ID:        "uuid",
					Name:      "Jamal",
					Email:     "jamalunyu@gmail.com",
					CreatedAt: time.Now(),
				},
			},
		},
		{
			Name: "Bad Request - Invalid JSON",
			Input: parameter{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://0.0.0.0/auth/register",
					strings.NewReader("invalid json"),
				),
			},
			mockBehavior: func(mockUsecase *MockIUserUsecase) {

			},
			expectedHeader: jsonHeader,
			expectedBody: response.HttpResponse{
				Status:  http.StatusBadRequest,
				Message: "",
				Data:    &dto.UserRegisterResponse{},
			},
		},
		{
			Name: "Conflict - Email already exists",
			Input: parameter{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://0.0.0.0/auth/register",
					strings.NewReader(`
					{
						"name":     "Jamal",
						"email":    "jamalunyu@gmail.com",
						"password": "Rahasia#123"
					},
				`),
				),
			},
			mockBehavior: func(mockUsecase *MockIUserUsecase) {
				mockUsecase.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(&dto.UserRegisterResponse{}, customErr.ErrEmailExist)
			},
			expectedHeader: jsonHeader,
			expectedBody: response.HttpResponse{
				Status:  http.StatusConflict,
				Message: "email already exist",
				Data:    &dto.UserRegisterResponse{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.mockBehavior(mockUsecase)

			userHandler.Register(tc.Input.w, tc.Input.r)

			rec := tc.Input.w.(*httptest.ResponseRecorder)
			res := rec.Result()

			if !reflect.DeepEqual(res.StatusCode, tc.expectedBody.Status) {
				t.Errorf("userHandler.Register() status code = %v, want %v", res.StatusCode, tc.expectedBody.Status)
			}

			if !reflect.DeepEqual(res.Header, tc.expectedHeader) {
				t.Errorf("userHandler.Register() header = %v, want %v", res.Header, tc.expectedHeader)
			}

			bodyBuffer := new(bytes.Buffer)
			bodyBuffer.ReadFrom(res.Body)

			var actualData response.HttpResponse
			json.Unmarshal(bodyBuffer.Bytes(), &actualData)

			expectedData := tc.expectedBody.Data.(*dto.UserRegisterResponse)

			if actualData.Status != http.StatusOK &&
				actualData.Status != http.StatusCreated {
				if actualData.Message == "" {
					t.Errorf("Error message should not be empty for failed response")
				}
				fmt.Println(actualData.Message)
				if actualData.Data != nil {
					t.Errorf("Data should be nil for error responses, got %v", actualData.Data)
				}

			} else {
				var actualDataBody *dto.UserRegisterResponse
				dataJSON, _ := json.Marshal(actualData.Data)
				json.Unmarshal(dataJSON, &actualDataBody)

				if expectedData.ID != actualDataBody.ID {
					t.Errorf("Id from response data = %s, want %s", actualDataBody.ID, expectedData.ID)
				}

				if expectedData.Name != actualDataBody.Name {
					t.Errorf("Name from response data = %s, want %s", actualDataBody.Name, expectedData.Name)
				}

				if expectedData.Email != actualDataBody.Email {
					t.Errorf("Email from response data = %s, want %s", actualDataBody.Email, expectedData.Email)
				}

				timeDiff := actualDataBody.CreatedAt.Sub(expectedData.CreatedAt)

				if timeDiff > time.Second || timeDiff < -time.Second {
					t.Errorf("Created at from response data = %v, want %v", actualDataBody.CreatedAt, expectedData.CreatedAt)
				}
			}
		})
	}
}
