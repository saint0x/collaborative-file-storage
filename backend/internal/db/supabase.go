package db

import (
	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	Client *supabase.Client
}

func NewSupabaseClient(supabaseURL, supabaseKey string) (*SupabaseClient, error) {
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, err
	}

	return &SupabaseClient{Client: client}, nil
}

// Add methods for database operations here
