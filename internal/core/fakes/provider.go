package fakes

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

func PipelineProvider(pipelines []sdk.Pipeline, updates sdk.PipelineUpdates) *Provider {
	return &Provider{
		Pipelines: pipelines,
		Updates:   updates,
	}
}

func AppProvider(apps []sdk.App) *Provider {
	return &Provider{
		Apps: apps,
	}
}

func RoutesProvider(routes map[string]sdk.AppRouting) *Provider {
	return &Provider{
		Routes: routes,
	}
}

func InstancesProvider(instances map[string]sdk.AppInstances) *Provider {
	return &Provider{
		Instances: instances,
	}
}

func NewErrProvider(err error) *Provider {
	return &Provider{
		Err: err,
	}
}

type Provider struct {
	Apps         []sdk.App
	Err          error
	Pipelines    []sdk.Pipeline
	Routes       map[string]sdk.AppRouting
	Instances    map[string]sdk.AppInstances
	Updates      sdk.PipelineUpdates
	RecordedTime time.Time
}

func (f *Provider) GetAppInstances(id string) (sdk.AppInstances, error) {
	return f.Instances[id], f.Err
}

func (f *Provider) GetAppRouting(id string) (sdk.AppRouting, error) {
	return f.Routes[id], f.Err
}

func (f *Provider) ListUpdates(since time.Time) (sdk.PipelineUpdates, error) {
	f.RecordedTime = since
	return f.Updates, nil
}

func (f Provider) ListPipelines() ([]sdk.Pipeline, error) {
	return f.Pipelines, nil
}

func (f Provider) GetPipeline(id string) (sdk.Pipeline, error) {
	panic("implement me")
}

func (f Provider) GetHistory(id string, before time.Time, limit int) (sdk.PipelineStatusList, error) {
	panic("implement me")
}

func (f Provider) ListApps() ([]sdk.App, error) {
	return f.Apps, f.Err
}

func (f Provider) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}
