package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ClerkService struct {
	BaseURL    string
	SecretKey  string
	HTTPClient *http.Client
}

type SessionClaims struct {
	Subject string `json:"sub"`
	// Add other relevant claims as needed
}

func NewClerkService() (*ClerkService, error) {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("CLERK_SECRET_KEY is not set")
	}

	return &ClerkService{
		BaseURL:   "https://api.clerk.com/v1",
		SecretKey: secretKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

// ValidateAndExtractUserID validates the token and extracts the user ID
func (cs *ClerkService) ValidateAndExtractUserID(ctx context.Context, token string) (string, error) {
	claims, err := cs.VerifyToken(token)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

// VerifyToken verifies the token and returns the session claims
func (cs *ClerkService) VerifyToken(token string) (*SessionClaims, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tokens/verify", cs.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cs.SecretKey))
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	resp, err := cs.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to verify token: %s", resp.Status)
	}

	var claims SessionClaims
	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// GetUser retrieves a user by ID
func (cs *ClerkService) GetUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", cs.BaseURL, userID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cs.SecretKey))

	resp, err := cs.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %s", resp.Status)
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers retrieves a list of users
func (cs *ClerkService) ListUsers(ctx context.Context) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", cs.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cs.SecretKey))

	resp, err := cs.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status)
	}

	var users []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}
