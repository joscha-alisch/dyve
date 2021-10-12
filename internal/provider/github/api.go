package github

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"
	"net/http"
)

/**
API is a simplified wrapper around the github api
*/
type API interface {
	ListTeams(org string) ([]Team, error)
	ListMembers(org string, team string) ([]Member, error)
}

type Login struct {
	AppId          int    `yaml:"appId"`
	InstallationId int    `yaml:"installationId"`
	PrivateKey     string `yaml:"privateKey"`
}

func NewDefaultApi(l Login) (API, error) {
	itr, err := ghinstallation.New(
		http.DefaultTransport,
		int64(l.AppId),
		int64(l.InstallationId),
		[]byte(l.PrivateKey),
	)
	if err != nil {
		return nil, err
	}

	cli := NewClient(&http.Client{Transport: itr})

	return NewApi(cli), nil
}

func NewApi(c Cli) API {
	return &api{
		c: c,
	}
}

type api struct {
	c Cli
}

func (a *api) ListMembers(org string, team string) ([]Member, error) {
	opt := &github.TeamListTeamMembersOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allUsers []*github.User
	for {
		users, resp, err := a.c.ListTeamMembersBySlug(context.Background(), org, team, opt)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var res []Member
	for _, user := range allUsers {
		name := user.GetName()
		if name == "" {
			name = user.GetLogin()
		}
		res = append(res, Member{
			Guid: fmt.Sprintf("%d", user.ID),
			Name: name,
		})
	}

	return res, nil
}

func (a *api) ListTeams(org string) ([]Team, error) {
	opt := &github.ListOptions{PerPage: 10}

	var allTeams []*github.Team
	for {
		teams, resp, err := a.c.ListTeams(context.Background(), org, opt)
		if err != nil {
			return nil, err
		}
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var res []Team
	for _, team := range allTeams {
		res = append(res, Team{
			TeamInfo: TeamInfo{
				Guid: fmt.Sprintf("%d", team.GetID()),
				Slug: team.GetSlug(),
				Name: team.GetName(),
			},
		})
	}

	return res, nil
}
