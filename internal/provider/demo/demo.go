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
	versions  map[string]sdk.PipelineVersionList
}

func (d *Provider) GenerateApps() {
	var apps []sdk.App
	for i := 0; i < 1000; i++ {
		namespace := namespace()
		apps = append(apps, sdk.App{
			Id:   id(),
			Name: appName(),
			Labels: sdk.AppLabels{
				"namespace": namespace,
				"version":   version(),
				"team":      team(),
			},
			Position: sdk.AppPosition{namespace},
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
	versions := make(map[string]sdk.PipelineVersionList)

	for i, pipeline := range d.pipelines {
		now := time.Now()
		start := now.Add(-(time.Hour * 24) * time.Duration(randomdata.Number(1, 365)))

		pVersion := sdk.PipelineVersion{
			PipelineId: pipeline.Id,
			Created:    start,
			Definition: d.generatePipelineDef(),
		}
		versions[pipeline.Id] = append(versions[pipeline.Id], pVersion)
		d.pipelines[i].Current = pVersion

		for start.Before(now) {
			run := sdk.PipelineStatus{
				Started:    start,
				PipelineId: pipeline.Id,
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
	d.versions = versions
}

func (d *Provider) generatePipelineDef() sdk.PipelineDefinition {
	def := sdk.PipelineDefinition{}
	hasEdge := make(map[int]int)
	id := 0
	columns := make([][]sdk.PipelineStep, randomdata.Number(2, 6))
	for col := 0; col < len(columns); col++ {
		rows := randomdata.Number(1, 4)
		for row := 0; row < rows; row++ {
			step := sdk.PipelineStep{
				Name:           pipelineStep(),
				Id:             id,
				AppDeployments: nil,
			}
			columns[col] = append(columns[col], step)
			def.Steps = append(def.Steps, step)

			conns := 0
			prevCol := col - 1
			if prevCol >= 0 {
				for _, prevStep := range columns[prevCol] {
					target := 0.2
					target += 0.05 * float64(hasEdge[prevStep.Id])
					target += float64(col-prevCol) * 0.2
					target += float64(conns) * 0.1
					if randomdata.Decimal(0, 1) > target {
						def.Connections = append(def.Connections, sdk.PipelineConnection{
							From:   prevStep.Id,
							To:     step.Id,
							Manual: randomdata.Boolean(),
						})
						hasEdge[prevStep.Id]++
						conns++
					}
				}
			}

			id++
		}
	}

	return def
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

func (d *Provider) GetHistory(id string, before time.Time, limit int) (sdk.PipelineStatusList, error) {
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

func (d *Provider) ListUpdates(since time.Time) (sdk.PipelineUpdates, error) {
	updates := sdk.PipelineUpdates{}

	for _, pipelineStatuses := range d.history {
		for _, i := range pipelineStatuses {
			if i.Started.After(since) {
				updates.Runs = append(updates.Runs, i)
				continue
			}

			for _, step := range i.Steps {
				if step.Started.After(since) || step.Ended.After(since) {
					updates.Runs = append(updates.Runs, i)
					break
				}
			}
		}
	}

	for _, pipelineVersions := range d.versions {
		for _, pipelineVersion := range pipelineVersions {
			if pipelineVersion.Created.After(since) {
				updates.Versions = append(updates.Versions, pipelineVersion)
			}
		}
	}

	return updates, nil
}
