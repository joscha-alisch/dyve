package client

import (
	"fmt"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
	"time"
)

func NewPipelineProviderClient(uri string, c *http.Client) sdk.PipelineProvider {
	return &pipelineProviderClient{
		baseClient: newBaseClient(uri+"/pipelines", c),
	}
}

type listPipelinesResponse struct {
	Status int
	Err    string
	Result []sdk.Pipeline
}

type getPipelineResponse struct {
	Status int
	Err    string
	Result sdk.Pipeline
}

type getHistoryResponse struct {
	Status int
	Err    string
	Result []sdk.PipelineStatus
}

type pipelineProviderClient struct {
	baseClient
}

func (p *pipelineProviderClient) ListPipelines() ([]sdk.Pipeline, error) {
	r := listPipelinesResponse{}
	err := p.get(&r, nil)
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}

func (p *pipelineProviderClient) GetPipeline(id string) (sdk.Pipeline, error) {
	r := getPipelineResponse{}
	err := p.get(&r, nil, id)
	if err != nil {
		return sdk.Pipeline{}, err
	}
	return r.Result, nil
}

func (p *pipelineProviderClient) GetHistory(id string, before time.Time, limit int) ([]sdk.PipelineStatus, error) {
	r := getHistoryResponse{}
	err := p.get(&r, map[string]string{
		"before": before.Format(time.RFC3339),
		"limit":  fmt.Sprintf("%d", limit),
	}, id, "history")
	if err != nil {
		return nil, err
	}
	return r.Result, nil
}
