package main

import (
	"github.com/joscha-alisch/dyve/internal/provider/cloudfoundry"
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

	cf, err := cloudfoundry.NewDefaultApi(cloudfoundry.Login{
		Api: c.cfApi, User: c.cfUser, Pass: c.cfPass,
	})
	if err != nil {
		panic(err)
	}

	db, err := cloudfoundry.NewMongoDatabase(c.mongoUri, c.mongoDb)
	if err != nil {
		panic(err)
	}

	r := cloudfoundry.NewReconciler(db, cf)
	s := cloudfoundry.NewScheduler(r)

	s.Run(8, 10*time.Second)

	<-make(chan struct{})
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
