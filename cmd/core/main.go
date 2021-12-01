package main

import (
	"flag"
	"fmt"
	"github.com/joscha-alisch/dyve/internal/core/api"
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/internal/core/config"
	coreDb "github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/groups"
	"github.com/joscha-alisch/dyve/internal/core/instances"
	"github.com/joscha-alisch/dyve/internal/core/pipelines"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	coreRecon "github.com/joscha-alisch/dyve/internal/core/reconciler"
	"github.com/joscha-alisch/dyve/internal/core/routing"
	"github.com/joscha-alisch/dyve/internal/core/service"
	"github.com/joscha-alisch/dyve/internal/core/teams"
	providerClient "github.com/joscha-alisch/dyve/internal/provider/client"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	c, err := config.LoadFrom(configPath)
	if err != nil {
		panic(err)
	}

	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	logLevel, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
		log.Error().Err(err).Msg("couldn't parse log level, setting to info")
	}
	zerolog.SetGlobalLevel(logLevel)

	db, err := coreDb.NewMongoDB(coreDb.MongoLogin{
		Uri: c.Database.URI,
		DB:  c.Database.Name,
	})
	if err != nil {
		panic(err)
	}

	providerService := provider.NewService(db)
	teamService := teams.NewService(db)
	appService := apps.NewService(db)
	groupService := groups.NewService(db, providerService)
	pipelineService := pipelines.NewService(db)
	routingService := routing.NewService(db)
	instancesService := instances.NewService(db)

	core := service.Core{
		Teams:     teamService,
		Apps:      appService,
		Groups:    groupService,
		Providers: providerService,
		Pipelines: pipelineService,
		Routing:   routingService,
		Instances: instancesService,
	}

	err = core.Pipelines.EnsureIndices()
	if err != nil {
		panic(err)
	}

	err = core.Teams.EnsureIndices()
	if err != nil {
		panic(err)
	}

	for _, providerConfig := range c.Providers {
		for _, feature := range providerConfig.Features {
			switch feature {
			case provider.TypeApps:
				p := providerClient.NewAppProviderClient(providerConfig.Host, nil)
				err = core.Providers.AddAppProvider(providerConfig.Id, providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypePipelines:
				p := providerClient.NewPipelineProviderClient(providerConfig.Host, nil)
				err = core.Providers.AddPipelineProvider(providerConfig.Id, providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypeGroups:
				p := providerClient.NewGroupProviderClient(providerConfig.Host, nil)
				err = core.Providers.AddGroupProvider(providerConfig.Id, providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypeRouting:
				p := providerClient.NewRoutingProviderClient(providerConfig.Host, nil)
				err = core.Providers.AddRoutingProvider(providerConfig.Id, providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypeInstances:
				p := providerClient.NewInstancesProviderClient(providerConfig.Host, nil)
				err = core.Providers.AddInstancesProvider(providerConfig.Id, providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	r := coreRecon.NewReconciler(core, time.Duration(c.Reconciliation.CacheSeconds)*time.Second)
	s := recon.NewScheduler(r)
	err = s.Run(8, 100*time.Millisecond)
	if err != nil {
		panic(err)
	}

	a := api.New(core, pipeviz.New(), api.Opts{
		DevConfig: c.DevConfig,
		Url:       c.ExternalUrl,
		Auth:      c.Auth,
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), a)
	if err != nil {
		panic(err)
	}
}
