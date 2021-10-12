package database

import (
	prov "github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

type Database interface {
	recon.JobProvider
	prov.Store

	ListAppsPaginated(perPage int, page int) (sdk.AppPage, error)
	GetApp(id string) (sdk.App, error)
	UpdateApps(providerId string, apps []sdk.App) error

	ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error)
	GetGroup(id string) (sdk.Group, error)
	UpdateGroups(providerId string, groups []sdk.Group) error

	ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error)
	GetPipeline(id string) (sdk.Pipeline, error)
	ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error)
	ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error)
	ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error)
	UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error
	AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error
	AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error
}

const (
	ReconcileAppProvider      recon.Type = "apps"
	ReconcilePipelineProvider recon.Type = "pipelines"
	ReconcileGroupProvider    recon.Type = "groups"
)
