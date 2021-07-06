package main

import (
	"github.com/joscha-alisch/dyve/internal/provider/cloudfoundry"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"os"
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
	c := getConfig()

	cf, err := cloudfoundry.NewDefaultApi(cloudfoundry.CFLogin{
		Api: c.cfApi, User: c.cfUser, Pass: c.cfPass,
	})
	if err != nil {
		panic(err)
	}

	db, err := cloudfoundry.NewMongoDatabase(cloudfoundry.MongoLogin{Uri: c.mongoUri, DB: c.mongoDb})
	if err != nil {
		panic(err)
	}

	r := cloudfoundry.NewReconciler(db, cf)
	s := cloudfoundry.NewScheduler(r)
	p := cloudfoundry.NewAppProvider(db)

	err = s.Run(8, 10*time.Second)
	if err != nil {
		panic(err)
	}

	err = sdk.ListenAndServeAppProvider(":9003", p)
	if err != nil {
		panic(err)
	}
}

func getConfig() config {
	return config{
		cfApi:    mustGetEnv("CF_API"),
		cfUser:   mustGetEnv("CF_USER"),
		cfPass:   mustGetEnv("CF_PASS"),
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
