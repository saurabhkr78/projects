package auth

import "context"

// contextKey is a custom type to avoid collisions with
// context keys from other packages.
type contextKey string

const authenticatedUserKey contextKey = "authenticated_user"

// IntoContext stores the authenticated user in the request context.
func IntoContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

// FromContext retrieves the authenticated user from the request context.
func FromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(authenticatedUserKey).(User)
	return user, ok
}
