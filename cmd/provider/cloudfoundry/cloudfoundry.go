package main

import (
	"github.com/joscha-alisch/dyve/internal/provider/cloudfoundry"
	"time"
)

func main() {
	cf := cloudfoundry.NewApi()
	db := cloudfoundry.NewDB()
	r := cloudfoundry.NewReconciler(db, cf)
	s := cloudfoundry.NewScheduler(r)

	s.Run(8, 10*time.Second)
}
