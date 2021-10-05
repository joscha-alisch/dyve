package api

import (
	"context"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/token"
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"net/http"
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
			service.AddProvider("github", opts.Auth.GitHub.Id, opts.Auth.GitHub.Secret)
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
