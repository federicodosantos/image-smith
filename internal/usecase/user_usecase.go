package usecase

import (
	"context"
	"time"

	"github.com/federicodosantos/image-smith/internal/dto"
	"github.com/federicodosantos/image-smith/internal/model"
	"github.com/federicodosantos/image-smith/internal/repository"
	customErr "github.com/federicodosantos/image-smith/pkg/error"
	"github.com/federicodosantos/image-smith/pkg/jwt"
	"github.com/federicodosantos/image-smith/pkg/regex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

type UserUsecase struct {
	userRepo repository.IUserRepository
	jwt      jwt.JWTItf
}

func NewUserUsecase(userRepo repository.IUserRepository, jwt jwt.JWTItf) IUserUsecase {
	return &UserUsecase{userRepo: userRepo, jwt: jwt}
}

func (u *UserUsecase) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, customErr.ErrEmailExist
	}

	if err := regex.Password(req.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createdUser := &model.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.userRepo.CreateUser(ctx, createdUser); err != nil {
		return nil, err
	}

	response := &dto.UserRegisterResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}

	return response, nil
}

func (u *UserUsecase) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := u.jwt.CreateToken(user.ID)

	return &dto.UserLoginResponse{
		JWTToken: token,
	}, nil
}
