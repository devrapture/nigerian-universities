package utils

import (
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	apperrors "github.com/coolpythoncodes/nigerian-universities/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJwt(userID uuid.UUID, email string, cfg *config.Config) (string, error) {
	expiryDuration := time.Duration(cfg.JwtExpires) * time.Hour
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JwtSecret))
}



func ValidateJwt(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}
	return claims, nil
}