package database

import (
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type Database interface {
	recon.JobProvider
	ListAppsPaginated(perPage int, page int) (sdk.AppPage, error)

	AddAppProvider(providerId string) error
	DeleteAppProvider(providerId string) error
	UpdateApps(providerId string, apps []sdk.App) error
}

const (
	ReconcileAppProvider recon.Type = "app-provider"
)
