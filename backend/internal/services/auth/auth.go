package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
)

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

// AuthenticateUser authenticates a user using the provided token
func AuthenticateUser(ctx context.Context, clerkService *ClerkService, token string) (context.Context, error) {
	userID, err := clerkService.ValidateAndExtractUserID(ctx, token)
	if err != nil {
		return ctx, err
	}

	return SetUserIDInContext(ctx, userID), nil
}

// Other auth-related functions...
