package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/federicodosantos/image-smith/internal/delivery"
	"github.com/federicodosantos/image-smith/internal/dto"
	response "github.com/federicodosantos/image-smith/pkg/response"
	"go.uber.org/mock/gomock"
)

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
					Return(&dto.UserRegisterResponse{}, nil)
			},
			expectedHeader: http.Header{
				"Content-Type": {"application/json"},
			},
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

			actualResponseBody := &dto.UserRegisterResponse{
				ID:        "uuid",
				Name:      "Jamal",
				Email:     "jamalunyu@gmail.com",
				CreatedAt: time.Now(),
			}
			json.Unmarshal(bodyBuffer.Bytes(), &actualResponseBody)

			expectedData := tc.expectedBody.Data.(*dto.UserRegisterResponse)

			if expectedData.ID != actualResponseBody.ID {
				t.Errorf("Id from response data = %s, want %s", actualResponseBody.ID, expectedData.ID)
			}

			if expectedData.Name != actualResponseBody.Name {
				t.Errorf("Name from response data = %s, want %s", actualResponseBody.Name, expectedData.Name)
			}

			if expectedData.Email != actualResponseBody.Email {
				t.Errorf("Email from response data = %s, want %s", actualResponseBody.Email, expectedData.Email)
			}

			timeDiff := actualResponseBody.CreatedAt.Sub(expectedData.CreatedAt)

			if timeDiff > time.Second || timeDiff < -time.Second {
				t.Errorf("Created at from response data = %v, want %v", actualResponseBody.CreatedAt, expectedData.CreatedAt)
			}

		})
	}
}
