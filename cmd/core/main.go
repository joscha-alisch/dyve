package main

import (
	"github.com/joscha-alisch/dyve/internal/core/api"
	coreDb "github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	coreRecon "github.com/joscha-alisch/dyve/internal/core/reconciler"
	providerClient "github.com/joscha-alisch/dyve/internal/provider/client"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"net/http"
	"os"
	"time"
)

func main() {
	c := getConfig()

	db, err := coreDb.NewMongoDB(coreDb.MongoLogin{
		Uri: c.mongoUri,
		DB:  c.mongoDb,
	})
	if err != nil {
		panic(err)
	}

	m := provider.NewManager(db)
	demoAppCli := providerClient.NewAppProviderClient("http://localhost:9003", nil)
	err = m.AddAppProvider("demo-apps", demoAppCli)
	if err != nil {
		panic(err)
	}

	demoPipelineCli := providerClient.NewPipelineProviderClient("http://localhost:9003", nil)
	err = m.AddPipelineProvider("demo-pipelines", demoPipelineCli)
	if err != nil {
		panic(err)
	}

	r := coreRecon.NewReconciler(db, m, 5*time.Second)
	s := recon.NewScheduler(r)
	err = s.Run(8, 5*time.Second)
	if err != nil {
		panic(err)
	}

	a := api.New(db)

	err = http.ListenAndServe(":9001", a)
	if err != nil {
		panic(err)
	}
}

type config struct {
	mongoUri string
	mongoDb  string
}

func getConfig() config {
	return config{
		mongoUri: mustGetEnv("MONGO_URI"),
		mongoDb:  mustGetEnv("MONGO_DB"),
	}
}

func mustGetEnv(env string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		panic("could not find env " + env)
	}
	return v
}
