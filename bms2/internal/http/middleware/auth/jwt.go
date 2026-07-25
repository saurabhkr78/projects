package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager is responsible for generating and validating JWTs.
type JWTManager struct {
	secret []byte
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
	}
}

// GenerateToken creates a signed JWT for the given user.
func (j *JWTManager) GenerateToken(user User, expiry time.Duration) (string, error) {
	claims := Claims{

		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

// ValidateToken parses, verifies and returns the JWT claims.
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {

			// Accept only HMAC signed tokens.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return j.secret, nil
		},
	)

	if err != nil {

		switch {

		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrExpiredToken

		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrMalformedToken

		default:
			return nil, ErrInvalidToken
		}
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims == nil {
		return nil, ErrMissingClaims
	}

	return claims, nil
}
