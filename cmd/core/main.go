package main

import (
	"flag"
	"fmt"
	"github.com/joscha-alisch/dyve/internal/core/api"
	"github.com/joscha-alisch/dyve/internal/core/config"
	coreDb "github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	coreRecon "github.com/joscha-alisch/dyve/internal/core/reconciler"
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

	m := provider.NewManager(db)
	for _, providerConfig := range c.Providers {
		for _, feature := range providerConfig.Features {
			switch feature {
			case provider.TypeApps:
				p := providerClient.NewAppProviderClient(providerConfig.Host, nil)
				err = m.AddAppProvider(providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypePipelines:
				p := providerClient.NewPipelineProviderClient(providerConfig.Host, nil)
				err = m.AddPipelineProvider(providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			case provider.TypeGroups:
				p := providerClient.NewGroupProviderClient(providerConfig.Host, nil)
				err = m.AddGroupProvider(providerConfig.Name, p)
				if err != nil {
					panic(err)
				}
			}
		}

	}

	r := coreRecon.NewReconciler(db, m, time.Duration(c.Reconciliation.CacheSeconds)*time.Second)
	s := recon.NewScheduler(r)
	err = s.Run(8, 10*time.Second)
	if err != nil {
		panic(err)
	}

	a := api.New(db, pipeviz.New(), api.Opts{
		DevMode: c.DevMode,
		Url:     c.ExternalUrl,
		Auth:    c.Auth,
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), a)
	if err != nil {
		panic(err)
	}
}
