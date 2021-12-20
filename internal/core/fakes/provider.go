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

func NewErrProvider(err error) *Provider {
	return &Provider{
		Err: err,
	}
}

type Provider struct {
	Apps         []sdk.App
	Err          error
	Pipelines    []sdk.Pipeline
	Updates      sdk.PipelineUpdates
	RecordedTime time.Time
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
