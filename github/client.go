package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client represents the wrapper of GitHub API client
type Client struct {
	client *github.Client
}

// NewClient creates new Client object
func NewClient(ctx context.Context, accesToken string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accesToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
	}
}
