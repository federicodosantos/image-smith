package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/federicodosantos/image-smith/internal/dto"
	"github.com/federicodosantos/image-smith/internal/usecase"
	customErr "github.com/federicodosantos/image-smith/pkg/error"
	response "github.com/federicodosantos/image-smith/pkg/response"
)

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func UserRoutes(router *http.ServeMux, userHandler *UserHandler) {
	router.HandleFunc("/auth/register", userHandler.Register)
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req *dto.UserRegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := uh.userUsecase.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, customErr.ErrEmailExist) {
			response.FailedResponse(w, http.StatusConflict, err.Error(), nil)
			return
		}
		response.FailedResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessResponse(w, http.StatusCreated, "successfully create user", user)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req *dto.UserLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := uh.userUsecase.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrEmailNotFound):
			response.FailedResponse(w, http.StatusNotFound, err.Error(), nil)
			return
		case errors.Is(err, customErr.ErrIncorrectPassword):
			response.FailedResponse(w, http.StatusUnauthorized, err.Error(), nil)
			return
		default:
			response.FailedResponse(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	response.SuccessResponse(w, http.StatusOK, "successfully login to account", token)
}
