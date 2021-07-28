package database

import (
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

type Database interface {
	recon.JobProvider
	ListAppsPaginated(perPage int, page int) (sdk.AppPage, error)
	GetApp(id string) (sdk.App, error)

	AddAppProvider(providerId string) error
	DeleteAppProvider(providerId string) error
	UpdateApps(providerId string, apps []sdk.App) error

	ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error)
	GetPipeline(id string) (sdk.Pipeline, error)
	ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error)

	AddPipelineProvider(providerId string) error
	DeletePipelineProvider(providerId string) error
	UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error
	AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error
}

const (
	ReconcileAppProvider      recon.Type = "apps"
	ReconcilePipelineProvider recon.Type = "pipelines"
)
