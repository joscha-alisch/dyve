package main

import (
	"flag"
	"fmt"
	"github.com/joscha-alisch/dyve/internal/provider/github"
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

	db, err := github.NewMongoDatabase(
		github.MongoLogin{
			Uri: c.Database.URI,
			DB:  c.Database.Name,
		}, c.GitHub.Org)

	gh, err := github.NewDefaultApi(c.GitHub.Login)
	if err != nil {
		panic(err)
	}

	r := github.NewReconciler(db, gh, 10*time.Minute)

	s := recon.NewScheduler(r)

	err = s.Run(8, 20*time.Second)
	if err != nil {
		panic(err)
	}

	p := github.NewGroupProvider(db)

	err = sdk.ListenAndServe(fmt.Sprintf(":%d", c.Port), sdk.ProviderConfig{
		Groups: p,
	})
	if err != nil {
		panic(err)
	}
}
