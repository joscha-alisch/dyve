package demo

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

func NewProvider() *Provider {
	p := &Provider{}
	p.GenerateApps()
	p.GeneratePipelines()
	p.GenerateHistory()

	return p
}

type Provider struct {
	apps      []sdk.App
	pipelines []sdk.Pipeline
	history   map[string][]sdk.PipelineStatus
}

func (d *Provider) GenerateApps() {
	var apps []sdk.App
	for i := 0; i < 1000; i++ {
		apps = append(apps, sdk.App{
			Id:   id(),
			Name: appName(),
			Meta: map[string]interface{}{
				"namespace": namespace(),
				"version":   version(),
				"team":      team(),
			},
		})
	}
	d.apps = apps
}

func (d *Provider) GeneratePipelines() {
	var pipelines []sdk.Pipeline
	for i := 0; i < 300; i++ {
		pipelines = append(pipelines, sdk.Pipeline{
			Id:   id(),
			Name: pipelineName(),
			Current: sdk.PipelineVersion{
				Definition: sdk.PipelineDefinition{
					Steps: []sdk.PipelineStep{
						{Name: "Build", Id: 0},
						{Name: "Deploy", Id: 1, AppDeployments: []string{"app"}},
					},
					Connections: []sdk.PipelineConnection{
						{
							From:   0,
							To:     1,
							Manual: false,
						},
					},
				},
			},
		})
	}
	d.pipelines = pipelines
}

func (d *Provider) GenerateHistory() {
	history := make(map[string][]sdk.PipelineStatus)
	for _, pipeline := range d.pipelines {
		now := time.Now()
		start := now.Add(-(time.Hour * 24) * time.Duration(randomdata.Number(1, 365)))

		for start.Before(now) {
			run := sdk.PipelineStatus{
				Started: start,
			}
			for _, step := range pipeline.Current.Definition.Steps {
				end := start.Add(time.Second * time.Duration(randomdata.Number(20, 180)))

				run.Steps = append(run.Steps, sdk.StepRun{
					StepId:  step.Id,
					Status:  sdk.StatusSuccess,
					Started: start,
					Ended:   end,
				})

				start = end.Add(time.Second * time.Duration(randomdata.Number(5, 20)))
			}
			history[pipeline.Id] = append(history[pipeline.Id], run)

			start = start.Add(time.Minute * time.Duration(randomdata.Number(30, 1000)))
		}
	}
	d.history = history
}

func (d *Provider) ListPipelines() ([]sdk.Pipeline, error) {
	return d.pipelines, nil
}

func (d *Provider) GetPipeline(id string) (sdk.Pipeline, error) {
	for _, p := range d.pipelines {
		if p.Id == id {
			return p, nil
		}
	}
	return sdk.Pipeline{}, sdk.ErrNotFound
}

func (d *Provider) GetHistory(id string, before time.Time, limit int) ([]sdk.PipelineStatus, error) {
	runs := d.history[id]

	var res []sdk.PipelineStatus
	for i := len(runs) - 1; i > 0; i-- {
		if runs[i].Started.Before(before) {
			res = append(res, runs[i])
			if len(res) == limit {
				break
			}
		}
	}

	return res, nil
}

func (d *Provider) ListApps() ([]sdk.App, error) {
	return d.apps, nil
}

func (d *Provider) GetApp(id string) (sdk.App, error) {
	for _, app := range d.apps {
		if app.Id == id {
			return app, nil
		}
	}
	return sdk.App{}, sdk.ErrNotFound
}
