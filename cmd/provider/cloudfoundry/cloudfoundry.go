package main

import (
	"fmt"
	"github.com/joscha-alisch/dyve/internal/provider/cloudfoundry"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

type config struct {
	cfApi    string
	cfUser   string
	cfPass   string
	mongoUri string
	mongoDb  string
}

func main() {
	c, err := LoadFrom("./config.yaml")
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
	p := cloudfoundry.NewAppProvider(db)

	err = s.Run(8, 10*time.Second)
	if err != nil {
		panic(err)
	}

	err = sdk.ListenAndServeAppProvider(fmt.Sprintf(":%d", c.Port), p)
	if err != nil {
		panic(err)
	}
}
