package github

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Client represents the wrapper of GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates new Client object
func NewClient(ctx context.Context, accessToken string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		ctx:    ctx,
	}
}

// CommitFronRef returns the latest commit SHA-1 of the given ref
// (branch, full commit SHA-1, short commit SHA-1...)
func (c *Client) CommitFronRef(repo, ref string) (string, error) {
	ss := strings.Split(repo, "/")
	if len(ss) != 2 {
		return "", errors.Errorf("invalid repository %q, must be owner/repo", repo)
	}

	sha1, _, err := c.client.Repositories.GetCommitSHA1(c.ctx, ss[0], ss[1], ref, "")
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve commit SHA-1")
	}

	return sha1, nil
}

// CreateDeployment creates Deployment and returns Deployment ID
// https://developer.github.com/v3/repos/deployments/
func (c *Client) CreateDeployment(repo, ref, cluster string) (int, error) {
	ss := strings.Split(repo, "/")
	if len(ss) != 2 {
		return -1, errors.Errorf("invalid repository %q, must be owner/repo", repo)
	}

	d, _, err := c.client.Repositories.CreateDeployment(c.ctx, ss[0], ss[1], &github.DeploymentRequest{
		Description:      github.String("k8ship deploy"),
		Environment:      github.String(cluster),
		Ref:              github.String(ref),
		RequiredContexts: &[]string{},
	})
	if err != nil {
		return -1, errors.Wrap(err, "failed to create Deployment")
	}

	_, _, err = c.client.Repositories.CreateDeploymentStatus(c.ctx, ss[0], ss[1], d.GetID(), &github.DeploymentStatusRequest{
		State: github.String("success"),
	})
	if err != nil {
		return -1, errors.Wrap(err, "failed to update Deployment status: 'pending' -> 'success'")
	}

	return d.GetID(), nil
}
