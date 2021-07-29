package api

import (
	"bytes"
	"encoding/json"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

func TestHttp(t *testing.T) {
	tests := []struct {
		desc   string
		state  fakeDb
		method string
		path   string
	}{
		{"gets app", fakeDb{app: sdk.App{
			Id: "guid-a", Name: "name-a",
		}}, "GET", "/api/apps/guid-a"},
		{"lists apps", fakeDb{appPage: sdk.AppPage{
			Pagination: sdk.Pagination{
				TotalResults: 20,
				TotalPages:   10,
				PerPage:      2,
				Page:         5,
			},
			Apps: []sdk.App{
				{Id: "guid-a", Name: "name-a"},
				{Id: "guid-b", Name: "name-b"},
			},
		}}, "GET", "/api/apps?perPage=2&page=5"},
		{"gets pipeline", fakeDb{pipeline: sdk.Pipeline{
			Id: "guid-a", Name: "name-a",
		}}, "GET", "/api/pipelines/guid-a"},
		{"lists pipelines", fakeDb{pipelinePage: sdk.PipelinePage{
			Pagination: sdk.Pagination{
				TotalResults: 20,
				TotalPages:   10,
				PerPage:      2,
				Page:         5,
			},
			Pipelines: []sdk.Pipeline{
				{Id: "guid-a", Name: "name-a"},
				{Id: "guid-b", Name: "name-b"},
			},
		}}, "GET", "/api/pipelines?perPage=2&page=5"},
		{"gets pipeline status", fakeDb{pipeline: sdk.Pipeline{
			Id:   "pipeline-a",
			Name: "pipeline",
			Current: sdk.PipelineVersion{
				Created: someTime.Add(-3 * time.Minute),
				Definition: sdk.PipelineDefinition{
					Steps: []sdk.PipelineStep{
						{
							Name:           "step-a",
							Id:             0,
							AppDeployments: nil,
						},
						{
							Name:           "step-b",
							Id:             1,
							AppDeployments: nil,
						},
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
		}, runs: []sdk.PipelineStatus{
			{
				PipelineId: "pipeline-a",
				Started:    someTime.Add(-2 * time.Minute),
				Steps: []sdk.StepRun{
					{
						StepId:  0,
						Status:  "succeeded",
						Started: someTime.Add(-2 * time.Minute),
						Ended:   someTime.Add(-1 * time.Minute),
					},
				},
			},
		}}, "GET", "/api/pipelines/pipeline-a/status"},
		{"gets pipeline runs", fakeDb{runs: []sdk.PipelineStatus{
			{
				PipelineId: "pipeline-a",
				Started:    someTime.Add(-2 * time.Minute),
				Steps: []sdk.StepRun{
					{
						StepId:  0,
						Status:  "succeeded",
						Started: someTime.Add(-2 * time.Minute),
						Ended:   someTime.Add(-1 * time.Minute),
					},
				},
			},
		}, versions: []sdk.PipelineVersion{
			{
				PipelineId: "pipeline-a",
				Created:    someTime.Add(-3 * time.Minute),
				Definition: sdk.PipelineDefinition{
					Steps: []sdk.PipelineStep{
						{
							Name:           "step-a",
							Id:             0,
							AppDeployments: nil,
						},
						{
							Name:           "step-b",
							Id:             1,
							AppDeployments: nil,
						},
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
		}}, "GET", "/api/pipelines/pipeline-a/runs"},
	}

	currentTime = func() time.Time {
		return someTime
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			h := New(&test.state, &test.state)

			s := httptest.NewServer(h)
			defer s.Close()

			w := httptest.NewRecorder()

			r := httptest.NewRequest(test.method, s.URL+test.path, nil)
			h.ServeHTTP(w, r)

			resp := w.Result()
			res, _ := httputil.DumpResponse(resp, false)
			m := make(map[string]interface{})
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			_ = json.Unmarshal(bodyBytes, &m)
			jsonString, _ := json.MarshalIndent(m, "", "    ")
			responseString := string(res) + string(jsonString)

			approvals.UseFolder("testdata")
			approvals.UseReporter(reporters.NewGoLandReporter())
			approvals.VerifyString(tt, responseString)
		})
	}

}

type fakeDb struct {
	appPage      sdk.AppPage
	app          sdk.App
	pipelinePage sdk.PipelinePage
	pipeline     sdk.Pipeline
	runs         []sdk.PipelineStatus
	versions     []sdk.PipelineVersion
}

func (f *fakeDb) ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error) {
	var res []sdk.PipelineStatus
	for _, run := range f.runs {
		if run.PipelineId == id && run.Started.Before(toExcl) {
			res = append(res, run)
		}
		if len(res) >= limit {
			break
		}
	}

	return res, nil
}

func (f *fakeDb) ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error) {
	var res sdk.PipelineVersionList
	for _, version := range f.versions {
		if version.PipelineId == id {
			res = append(res, version)
		}
	}
	return res, nil
}

func (f *fakeDb) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	panic("implement me")
}

func (f *fakeDb) ListPipelineRuns(id string, from time.Time, to time.Time) (sdk.PipelineStatusList, error) {
	var res []sdk.PipelineStatus
	for _, run := range f.runs {
		if run.PipelineId == id && run.Started.After(from) && run.Started.Before(to) {
			res = append(res, run)
		}
	}

	return res, nil
}

func (f *fakeDb) Write(graph pipeviz.Graph, w io.Writer) {
	b, _ := json.Marshal(graph)
	_, _ = w.Write(b)
}

func (f *fakeDb) Generate(graph pipeviz.Graph) []byte {
	var buf bytes.Buffer
	buf.Write([]byte("fake svg: "))
	f.Write(graph, &buf)
	return buf.Bytes()
}

func (f *fakeDb) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	if f.pipelinePage.PerPage == perPage && f.pipelinePage.Page == page {
		return f.pipelinePage, nil
	}
	return sdk.PipelinePage{}, sdk.ErrPageExceeded
}

func (f *fakeDb) GetPipeline(id string) (sdk.Pipeline, error) {
	if f.pipeline.Id == id {
		return f.pipeline, nil
	}
	return sdk.Pipeline{}, sdk.ErrNotFound
}

func (f *fakeDb) AddPipelineProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) DeletePipelineProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	panic("implement me")
}

func (f *fakeDb) GetApp(id string) (sdk.App, error) {
	if f.app.Id == id {
		return f.app, nil
	}
	return sdk.App{}, sdk.ErrNotFound
}

func (f *fakeDb) AddAppProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) DeleteAppProvider(providerId string) error {
	panic("implement me")
}

func (f *fakeDb) UpdateApps(providerId string, apps []sdk.App) error {
	panic("implement me")
}

func (f *fakeDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	panic("implement me")
}

func (f *fakeDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	if f.appPage.PerPage == perPage && f.appPage.Page == page {
		return f.appPage, nil
	}
	return sdk.AppPage{}, sdk.ErrPageExceeded
}
