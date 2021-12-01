package main

import (
	"flag"
	"fmt"
	"github.com/joscha-alisch/dyve/internal/provider/cloudfoundry"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
}

func main() {
	flag.Parse()
	c, err := LoadFrom(configPath)
	if err != nil {
		panic(err)
	}

	cf, err := cloudfoundry.NewDefaultApi(cloudfoundry.CFLogin{
		Api: c.CloudFoundry.Api, User: c.CloudFoundry.User, Pass: c.CloudFoundry.Password,
	})
	if err != nil {
		panic(err)
	}

	db, err := cloudfoundry.NewMongoDatabase(cloudfoundry.MongoLogin{Uri: c.Database.URI, DB: c.Database.Name})
	if err != nil {
		panic(err)
	}

	r := cloudfoundry.NewReconciler(db, cf, time.Duration(c.Reconciliation.CacheSeconds)*time.Second)
	s := recon.NewScheduler(r)
	p := cloudfoundry.NewProvider(db, cf)

	err = s.Run(8, 10*time.Second)
	if err != nil {
		panic(err)
	}

	err = sdk.ListenAndServe(fmt.Sprintf(":%d", c.Port), sdk.ProviderConfig{
		Apps:      p,
		Routing:   p,
		Instances: p,
	})
	if err != nil {
		panic(err)
	}
}
