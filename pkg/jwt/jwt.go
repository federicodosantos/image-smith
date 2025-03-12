package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTItf interface {
	CreateToken(userID string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type JWT struct {
	SecretKey  string
	ExpireTime time.Duration
}

func NewJwt(SecretKey string, ExpireTime string) (JWTItf, error) {
	exp, err := time.ParseDuration(ExpireTime)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format for expireTime: %v", err)
	}

	return &JWT{
		SecretKey:  SecretKey,
		ExpireTime: exp,
	}, nil
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserID string
}

// CreateToken implements JWTItf.
func (j *JWT) CreateToken(userID string) (string, error) {
	if j.ExpireTime <= 0 {
		return "", fmt.Errorf("jwt expire time must be greater than 0")
	}

	claims := &UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// VerifyToken implements JWTItf.
func (j *JWT) VerifyToken(tokenString string) (string, error) {
	var claims UserClaim

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
