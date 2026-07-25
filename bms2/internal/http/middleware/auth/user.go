package auth

// User represents the authenticated user extracted from the JWT.
//
// This is the identity that the rest of the application works with.
// It is intentionally kept separate from JWT Claims.
type User struct {
	ID    string
	Email string
	Role  string
}
