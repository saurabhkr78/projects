package auth

import "github.com/golang-jwt/jwt/v5"

// Claims represents the data stored inside the JWT payload.
//
// Along with our custom fields (UserID, Email, Role),
// we embed jwt.RegisteredClaims to get standard JWT fields
// like exp (expiry), iat (issued at), nbf (not before), iss (issuer), etc.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}
