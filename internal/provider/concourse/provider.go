package concourse

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

func NewProvider() sdk.PipelineProvider {
	return &provider{}
}

type provider struct {
}

func (p *provider) ListPipelines() ([]sdk.Pipeline, error) {
	panic("implement me")
}

func (p *provider) ListUpdates(since time.Time) (sdk.PipelineUpdates, error) {
	panic("implement me")
}

func (p *provider) GetPipeline(id string) (sdk.Pipeline, error) {
	panic("implement me")
}

func (p *provider) GetHistory(id string, before time.Time, limit int) (sdk.PipelineStatusList, error) {
	panic("implement me")
}
