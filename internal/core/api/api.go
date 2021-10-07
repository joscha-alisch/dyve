package api

import (
	"context"
	"fmt"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/token"
	"github.com/google/go-github/v39/github"
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

type Opts struct {
	Url     string
	DevMode bool
	Auth    config.AuthConfig
}

func New(db database.Database, pipeGen pipeviz.PipeViz, opts Opts) http.Handler {
	a := &api{
		Router:  mux.NewRouter(),
		db:      db,
		pipeGen: pipeGen,
	}

	if opts.Auth.Secret == "" {
		panic("Need to provide an auth secret")
	}

	authOpts := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return opts.Auth.Secret, nil
		}),
		TokenDuration:   time.Minute * 5,
		CookieDuration:  time.Hour * 24,
		Issuer:          "dyve",
		URL:             opts.Url,
		AvatarStore:     avatar.NewLocalFS("/tmp"),
		AvatarRoutePath: "/auth/avatars",
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			if opts.Auth.GitHub.Enabled && strings.HasPrefix(claims.User.ID, "github") {
				if !userIsInOrg(claims.User, opts.Auth.GitHub.Org) {
					log.Debug().
						Str("user", claims.User.Name).
						Str("required", opts.Auth.GitHub.Org).
						Strs("orgs", claims.User.SliceAttr("orgs")).
						Msg("token declined because user is not in org")
					return false
				}
			}

			return claims.User != nil
		}),
	}

	// create auth service with providers
	service := auth.NewService(authOpts)
	if opts.DevMode {
		service.AddDevProvider(8000)

		go func() {
			devAuthServer, err := service.DevAuth()
			if err != nil {
				panic(err)
			}
			devAuthServer.Run(context.Background())
		}()
	} else {
		if opts.Auth.GitHub.Enabled {
			service.AddProviderWithOptions("github", opts.Auth.GitHub.Id, opts.Auth.GitHub.Secret, []string{"read:org"}, func(c *http.Client, u token.User) token.User {
				gh := github.NewClient(c)
				t, _, _ := gh.Teams.ListUserTeams(context.Background(), &github.ListOptions{})

				orgs := make(map[string]bool)
				var teams []string
				for _, team := range t {
					orgs[team.Organization.GetLogin()] = true
					teams = append(teams, fmt.Sprintf("%s:%s", team.Organization.GetLogin(), team.GetName()))
				}

				var orgList []string
				for org := range orgs {
					orgList = append(orgList, org)
				}

				u.SetSliceAttr("orgs", orgList)
				u.SetSliceAttr("teams", teams)

				log.Debug().Str("user", u.Name).Msg("new login")
				return u
			})
		}
	}

	authRoutes, avaRoutes := service.Handlers()
	a.PathPrefix("/auth/avatars").Handler(avaRoutes)
	a.PathPrefix("/auth").Handler(authRoutes)

	authenticated := service.Middleware()
	api := a.PathPrefix("/api").Subrouter()
	api.Use(authenticated.Auth)

	api.Path("/apps").Queries("perPage", "").HandlerFunc(a.listAppsPaginated)
	api.Path("/apps/{id:[0-9a-z-]+}").HandlerFunc(a.getApp)

	api.Path("/pipelines").Queries("perPage", "").HandlerFunc(a.listPipelinesPaginated)
	api.Path("/pipelines/{id:[0-9a-z-]+}/status").HandlerFunc(a.getPipelineStatus)
	api.Path("/pipelines/{id:[0-9a-z-]+}/runs").HandlerFunc(a.listPipelineRuns)
	api.Path("/pipelines/{id:[0-9a-z-]+}").HandlerFunc(a.getPipeline)

	return a
}

type api struct {
	*mux.Router
	db      database.Database
	pipeGen pipeviz.PipeViz
}

func userIsInOrg(user *token.User, org string) bool {
	orgs := user.SliceAttr("orgs")
	for _, s := range orgs {
		if s == org {
			return true
		}
	}
	return false
}
