package api

import (
	"context"
	"errors"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/live"
	"github.com/joscha-alisch/dyve/internal/core/service"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"net/http"
	"time"
)

type Opts struct {
	Url       string
	DevConfig config.DevConfig
	Auth      config.AuthConfig
}

func New(core service.Core, pipeGen pipeviz.PipeViz, opts Opts) http.Handler {
	a := &api{
		Router:             mux.NewRouter(),
		core:               core,
		pipeGen:            pipeGen,
		appViewer:          live.NewAppViewer(core),
		disableOriginCheck: opts.DevConfig.DisableOriginCheck,
	}

	if opts.Auth.Secret == "" && opts.DevConfig.DisableAuth == false {
		panic("Need to provide an auth secret")
	}

	authOpts := auth.Opts{
		SecretReader:    getTokenSecretFunc(opts),
		TokenDuration:   time.Minute * 5,
		CookieDuration:  time.Hour * 24,
		Issuer:          "dyve",
		URL:             opts.Url,
		AvatarStore:     avatar.NewLocalFS("/tmp"),
		AvatarRoutePath: "/auth/avatars",
		ClaimsUpd:       getUpdateClaimsFunc(opts),
		Validator:       getTokenValidatorFunc(opts),
	}

	// create auth service with providers
	service := auth.NewService(authOpts)
	if opts.DevConfig.UseFakeOauth2 {
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
			service.AddProviderWithOptions("github", opts.Auth.GitHub.Id, opts.Auth.GitHub.Secret, []string{"read:org"}, getGHProviderFunc())
		}
	}

	authRoutes, avaRoutes := service.Handlers()
	a.PathPrefix("/auth/avatars").Handler(avaRoutes)
	a.PathPrefix("/auth").Handler(authRoutes)

	authenticated := service.Middleware()
	api := a.PathPrefix("/api").Subrouter()

	if !opts.DevConfig.DisableAuth {
		api.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Upgrade") == "websocket" {
					c, err := r.Cookie("XSRF-TOKEN")
					if err != nil {
						respondErr(w, http.StatusForbidden, errors.New("XSRF-TOKEN cookie not set"))
						return
					}
					r.Header.Set("X-XSRF-TOKEN", c.Value)
				}
				next.ServeHTTP(w, r)
			})
		})
		api.Use(authenticated.Auth)
	}

	api.Path("/apps").Queries("perPage", "").Methods("GET").HandlerFunc(a.listAppsPaginated)
	api.Path("/apps/{id:[0-9a-z-]+}/live").HandlerFunc(a.startWebsocketApp)
	api.Path("/apps/{id:[0-9a-z-]+}").Methods("GET").HandlerFunc(a.getApp)

	api.Path("/pipelines").Queries("perPage", "").HandlerFunc(a.listPipelinesPaginated)
	api.Path("/pipelines/{id:[0-9a-z-]+}/status").HandlerFunc(a.getPipelineStatus)
	api.Path("/pipelines/{id:[0-9a-z-]+}/runs").HandlerFunc(a.listPipelineRuns)
	api.Path("/pipelines/{id:[0-9a-z-]+}").HandlerFunc(a.getPipeline)

	api.Path("/teams").Queries("perPage", "").HandlerFunc(a.listTeamsPaginated)
	api.Path("/teams/{id:[0-9a-z-]+}").Methods("GET").HandlerFunc(a.getTeam)
	api.Path("/teams/{id:[0-9a-z-]+}").Methods("DELETE").HandlerFunc(a.deleteTeam)
	api.Path("/teams/{id:[0-9a-z-]+}").Methods("POST").HandlerFunc(a.createTeam)
	api.Path("/teams/{id:[0-9a-z-]+}").Methods("PUT").HandlerFunc(a.updateTeam)

	api.Path("/groups").HandlerFunc(a.listGroups)

	a.appViewer.Run()

	return a
}

type api struct {
	*mux.Router
	db                 database.Database
	pipeGen            pipeviz.PipeViz
	core               service.Core
	disableOriginCheck bool
	appViewer          *live.AppViewer
}
