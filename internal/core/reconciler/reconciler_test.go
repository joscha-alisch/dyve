package reconciler

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	tests := []struct {
		desc           string
		job            recon.Job
		before         map[string][]interface{}
		providerId     string
		provider       fakeProvider
		after          map[string][]interface{}
		expectedErr    error
		expectedWorked bool
	}{
		{desc: "adds apps", job: recon.Job{
			Type: database.ReconcileAppProvider,
			Guid: "app-provider",
		}, providerId: "app-provider", provider: fakeProvider{apps: []sdk.App{
			{Id: "app-a", Name: "app-a"},
			{Id: "app-b", Name: "app-b"},
		}}, after: map[string][]interface{}{
			"app-provider": {
				&sdk.App{Id: "app-a", Name: "app-a"},
				&sdk.App{Id: "app-b", Name: "app-b"},
			},
		}, expectedWorked: true},
		{desc: "adds pipelines", job: recon.Job{
			Type: database.ReconcilePipelineProvider,
			Guid: "pipeline-provider",
		}, providerId: "pipeline-provider", provider: fakeProvider{pipelines: []sdk.Pipeline{
			{Id: "pipeline-a", Name: "pipeline-a"},
			{Id: "pipeline-b", Name: "pipeline-b"},
		}}, after: map[string][]interface{}{
			"pipeline-provider": {
				&sdk.Pipeline{Id: "pipeline-a", Name: "pipeline-a"},
				&sdk.Pipeline{Id: "pipeline-b", Name: "pipeline-b"},
			},
		}, expectedWorked: true},
		{desc: "removes apps if provider not found", job: recon.Job{
			Type: database.ReconcileAppProvider,
			Guid: "not-exist",
		}, before: map[string][]interface{}{
			"not-exist": {
				sdk.App{Id: "app-a", Name: "app-a"},
				sdk.App{Id: "app-b", Name: "app-b"},
			},
		}, after: map[string][]interface{}{}, expectedWorked: true},
		{desc: "removes pipelines if provider not found", job: recon.Job{
			Type: database.ReconcilePipelineProvider,
			Guid: "not-exist",
		}, before: map[string][]interface{}{
			"not-exist": {
				sdk.Pipeline{Id: "pipeline-a", Name: "pipeline-a"},
				sdk.Pipeline{Id: "pipeline-b", Name: "pipeline-b"},
			},
		}, after: map[string][]interface{}{}, expectedWorked: true},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			db := &fakeDb{job: test.job, content: test.before}
			r := NewReconciler(db, &fakeManager{
				test.providerId,
				&test.provider,
			}, 1*time.Minute)
			worked, err := r.Run()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err %v\n   got err %v", test.expectedErr, err)
			}

			if worked != test.expectedWorked {
				tt.Errorf("\nwanted worked: %v\n   got worked: %v", test.expectedWorked, worked)
			}

			if !cmp.Equal(test.after, db.content) {
				tt.Errorf("\nstate diff: \n%s\n", cmp.Diff(test.after, db.content))
			}
		})
	}

}

type fakeManager struct {
	providerId string
	provider   *fakeProvider
}

func (f *fakeManager) AddPipelineProvider(id string, p sdk.PipelineProvider) error {
	panic("implement me")
}

func (f *fakeManager) GetPipelineProvider(id string) (sdk.PipelineProvider, error) {
	if id == f.providerId {
		return f.provider, nil
	}
	return nil, provider.ErrNotFound
}

func (f *fakeManager) AddAppProvider(id string, p sdk.AppProvider) error {
	panic("implement me")
}

func (f *fakeManager) GetAppProvider(id string) (sdk.AppProvider, error) {
	if id == f.providerId {
		return f.provider, nil
	}
	return nil, provider.ErrNotFound
}

type fakeProvider struct {
	apps      []sdk.App
	err       error
	pipelines []sdk.Pipeline
}

func (f fakeProvider) ListPipelines() ([]sdk.Pipeline, error) {
	return f.pipelines, nil
}

func (f fakeProvider) GetPipeline(id string) (sdk.Pipeline, error) {
	panic("implement me")
}

func (f fakeProvider) GetHistory(id string, before time.Time, limit int) ([]sdk.PipelineStatus, error) {
	panic("implement me")
}

func (f fakeProvider) ListApps() ([]sdk.App, error) {
	return f.apps, f.err
}

func (f fakeProvider) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}

type fakeDb struct {
	job     recon.Job
	content map[string][]interface{}
}

func (f *fakeDb) ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error) {
	panic("implement me")
}

func (f *fakeDb) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	panic("implement me")
}

func (f *fakeDb) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	panic("implement me")
}

func (f *fakeDb) GetPipeline(id string) (sdk.Pipeline, error) {
	panic("implement me")
}

func (f *fakeDb) AddPipelineProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) DeletePipelineProvider(providerId string) error {
	delete(f.content, providerId)
	return nil
}

func (f *fakeDb) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	if f.content == nil {
		f.content = make(map[string][]interface{})
	}

	f.content[providerId] = make([]interface{}, len(pipelines))
	for i, p := range pipelines {
		p := p
		f.content[providerId][i] = &p
	}
	return nil
}

func (f *fakeDb) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}

func (f *fakeDb) AddAppProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) DeleteAppProvider(providerId string) error {
	delete(f.content, providerId)
	return nil
}

func (f *fakeDb) UpdateApps(providerId string, apps []sdk.App) error {
	if f.content == nil {
		f.content = make(map[string][]interface{})
	}

	f.content[providerId] = make([]interface{}, len(apps))
	for i, app := range apps {
		app := app
		f.content[providerId][i] = &app
	}
	return nil
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	return f.job, true
}

func (f *fakeDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	panic("implement me")
}
