package api

import (
	"context"
	"fmt"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func getUpdateClaimsFunc(opts Opts) token.ClaimsUpdFunc {
	return func(claims token.Claims) token.Claims {
		if opts.DevConfig.UseFakeOauth2 && claims.User != nil {
			claims.User.SetSliceAttr("groups", opts.DevConfig.UserGroups)
		}
		return claims
	}
}

func getTokenValidatorFunc(opts Opts) token.ValidatorFunc {
	return func(_ string, claims token.Claims) bool {
		if opts.Auth.GitHub.Enabled && strings.HasPrefix(claims.User.ID, "github") {
			if !userIsInOrg(claims.User, opts.Auth.GitHub.Org) {
				log.Debug().
					Str("user", claims.User.Name).
					Str("required", opts.Auth.GitHub.Org).
					Strs("orgs", getUserOrgs(claims.User)).
					Msg("token declined because user is not in org")
				return false
			}
		}

		return claims.User != nil
	}
}

func getTokenSecretFunc(opts Opts) token.SecretFunc {
	return func(id string) (string, error) {
		return opts.Auth.Secret, nil
	}
}

func getGHProviderFunc() provider.ExtraUserInfoFunc {
	return func(c *http.Client, u token.User) token.User {
		gh := github.NewClient(c)
		t, _, _ := gh.Teams.ListUserTeams(context.Background(), &github.ListOptions{})

		orgs := make(map[string]bool)
		var teams []string
		for _, team := range t {
			orgs[team.Organization.GetLogin()] = true
			teams = append(teams, fmt.Sprintf("%s:%s:%d", "github", team.Organization.GetLogin(), team.GetID()))
		}

		var orgList []string
		for org := range orgs {
			orgList = append(orgList, org)
		}

		u.SetSliceAttr("orgs", orgList)
		u.SetSliceAttr("groups", teams)

		log.Debug().Str("user", u.Name).Msg("new login")
		return u
	}
}

func userIsInOrg(user *token.User, org string) bool {
	orgs := getUserOrgs(user)
	for _, s := range orgs {
		if s == org {
			return true
		}
	}
	return false
}

func getUserOrgs(u *token.User) []string {
	orgs, ok := u.Attributes["orgs"].([]interface{})
	if !ok {
		return nil
	}

	var res []string
	for _, org := range orgs {
		if orgString, ok := org.(string); ok {
			res = append(res, orgString)
		}
	}

	return res
}
