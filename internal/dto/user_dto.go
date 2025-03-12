package dto

import "time"

type UserRegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type UserLoginRequest struct {
	Email    string
	Password string
}

type UserRegisterResponse struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserLoginResponse struct {
	JWTToken string
}
