package api

import (
	"encoding/json"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"github.com/joscha-alisch/dyve/internal/core/fakes"
	"github.com/joscha-alisch/dyve/internal/core/service"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

func TestHttp(t *testing.T) {
	tests := []struct {
		desc              string
		state             state
		method            string
		path              string
		apps              *fakes.RecordingAppsService
		expectedApps      *fakes.AppsRecorder
		pipelines         *fakes.RecordingPipelinesService
		expectedPipelines *fakes.PipelinesRecorder
	}{
		{desc: "gets app", method: "GET", path: "/api/apps/guid-a", apps: &fakes.RecordingAppsService{
			App: apps.App{
				ProviderId: "provider",
				App: sdk.App{
					Id: "guid-a", Name: "name-a", Labels: map[string]string{"key": "value"}, Position: []string{"position"},
				},
			},
		}, expectedApps: &fakes.AppsRecorder{
			AppId: "guid-a",
		}},
		{desc: "lists apps", method: "GET", path: "/api/apps?perPage=2&page=5", apps: &fakes.RecordingAppsService{
			Page: sdk.AppPage{
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
			}}, expectedApps: &fakes.AppsRecorder{
			PerPage: 2,
			Page:    5,
		}},
		{desc: "gets pipeline", method: "GET", path: "/api/pipelines/guid-a", pipelines: &fakes.RecordingPipelinesService{
			Pipeline: sdk.Pipeline{
				Id: "guid-a", Name: "name-a",
			},
		}, expectedPipelines: &fakes.PipelinesRecorder{PipelineId: "guid-a"}},
		{desc: "lists pipelines", method: "GET", path: "/api/pipelines?perPage=2&page=5", pipelines: &fakes.RecordingPipelinesService{
			Page: sdk.PipelinePage{
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
			}}, expectedPipelines: &fakes.PipelinesRecorder{
			PerPage: 2,
			Page:    5,
		}},
		{desc: "gets pipeline status", method: "GET", path: "/api/pipelines/pipeline-a/status", pipelines: &fakes.RecordingPipelinesService{
			Pipeline: sdk.Pipeline{
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
			}, Runs: []sdk.PipelineStatus{
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
			},
		}, expectedPipelines: &fakes.PipelinesRecorder{
			PipelineId: "pipeline-a",
			FromIncl:   someTime.Add(-3 * time.Minute),
			ToExcl:     someTime,
		}},
		{desc: "gets pipeline runs", method: "GET", path: "/api/pipelines/pipeline-a/runs", pipelines: &fakes.RecordingPipelinesService{Runs: []sdk.PipelineStatus{
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
		}, Versions: []sdk.PipelineVersion{
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
		}}, expectedPipelines: &fakes.PipelinesRecorder{
			PipelineId: "pipeline-a",
			FromIncl:   someTime.Add(-2 * time.Minute),
			ToExcl:     someTime,
			Limit:      10,
		}},
		{desc: "gets empty pipeline runs", method: "GET", path: "/api/pipelines/pipeline-a/runs",
			pipelines: &fakes.RecordingPipelinesService{Runs: []sdk.PipelineStatus{}, Versions: []sdk.PipelineVersion{
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
			}}, expectedPipelines: &fakes.PipelinesRecorder{
				PipelineId: "pipeline-a",
				ToExcl:     someTime,
				Limit:      10,
			}},
	}

	currentTime = func() time.Time {
		return someTime
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			h := New(service.Core{
				Apps:      test.apps,
				Pipelines: test.pipelines,
			}, &fakes.PipeViz{}, Opts{
				DevConfig: config.DevConfig{DisableAuth: true},
			})

			testHttp(tt, h, test.method, test.path)

			if test.expectedPipelines != nil && !cmp.Equal(*test.expectedPipelines, test.pipelines.Record) {
				tt.Errorf("pipeline records don't match:%s\n", cmp.Diff(*test.expectedPipelines, test.pipelines.Record))
			}

			if test.expectedApps != nil && !cmp.Equal(*test.expectedApps, test.apps.Record) {
				tt.Errorf("app records don't match:%s\n", cmp.Diff(*test.expectedApps, test.apps.Record))
			}
		})
	}

}

func testHttp(tt *testing.T, h http.Handler, method string, path string) {
	s := httptest.NewServer(h)
	defer s.Close()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, s.URL+path, nil)
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
}

type state struct {
	appPage      sdk.AppPage
	app          sdk.App
	pipelinePage sdk.PipelinePage
	pipeline     sdk.Pipeline
	runs         []sdk.PipelineStatus
	versions     []sdk.PipelineVersion
}
