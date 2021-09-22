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

	db, err := coreDb.NewMongoDB(coreDb.MongoLogin{
		Uri: c.Database.URI,
		DB:  c.Database.Database,
	})
	if err != nil {
		panic(err)
	}

	m := provider.NewManager(db)
	for _, appProvider := range c.AppProviders {
		p := providerClient.NewAppProviderClient(appProvider.Host, nil)
		err = m.AddAppProvider(appProvider.Name, p)
		if err != nil {
			panic(err)
		}
	}

	r := coreRecon.NewReconciler(db, m, 5*time.Second)
	s := recon.NewScheduler(r)
	err = s.Run(8, 5*time.Second)
	if err != nil {
		panic(err)
	}

	a := api.New(db, pipeviz.New())

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), a)
	if err != nil {
		panic(err)
	}
}
