package auth

import "errors"

var (
	// ErrInvalidToken indicates the token failed validation.
	// Example: invalid signature, unsupported signing method, etc.
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken indicates the token's expiration time has passed.
	ErrExpiredToken = errors.New("token expired")

	// ErrMalformedToken indicates the JWT format is invalid.
	// Example: missing parts, invalid base64 encoding, etc.
	ErrMalformedToken = errors.New("malformed token")

	// ErrMissingClaims indicates the expected claims could not be extracted.
	ErrMissingClaims = errors.New("missing claims")
)
