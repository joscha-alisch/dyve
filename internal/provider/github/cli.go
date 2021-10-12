package github

import (
	"context"
	"github.com/google/go-github/v39/github"
	"net/http"
)

/**
Cli is an interface wrapping the methods we need from the official client
*/
type Cli interface {
	ListTeams(ctx context.Context, org string, opts *github.ListOptions) ([]*github.Team, *github.Response, error)
	ListTeamMembersBySlug(ctx context.Context, org string, slug string, opts *github.TeamListTeamMembersOptions) ([]*github.User, *github.Response, error)
}

func NewClient(c *http.Client) Cli {
	if c == nil {
		c = http.DefaultClient
	}

	ghCli := github.NewClient(c)

	return gitHubCli{
		ghCli.Teams,
	}
}

type gitHubCli struct {
	*github.TeamsService
}
