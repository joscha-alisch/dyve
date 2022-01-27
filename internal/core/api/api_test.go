package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/internal/core/config"
	"github.com/joscha-alisch/dyve/internal/core/fakes"
	"github.com/joscha-alisch/dyve/internal/core/fakes/fakeGroups"
	"github.com/joscha-alisch/dyve/internal/core/groups"
	"github.com/joscha-alisch/dyve/internal/core/service"
	"github.com/joscha-alisch/dyve/internal/core/teams"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

var someErr = errors.New("some error")

func TestHttp(t *testing.T) {
	tests := []struct {
		desc              string
		method            string
		path              string
		apps              *fakes.RecordingAppsService
		expectedApps      *fakes.AppsRecorder
		pipelines         *fakes.RecordingPipelinesService
		expectedPipelines *fakes.PipelinesRecorder
		teams             *fakes.RecordingTeamsService
		expectedTeams     *fakes.TeamsRecorder
		body              string
		groups            *fakeGroups.RecordingGroupsService
		expectedGroups    *fakeGroups.GroupsRecorder
		headers           http.Header
		overrideRequest   *http.Request
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
		{desc: "error while getting app", method: "GET", path: "/api/apps/guid-a", apps: &fakes.RecordingAppsService{
			Err: someErr,
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
		{desc: "lists apps with page default", method: "GET", path: "/api/apps?perPage=2", apps: &fakes.RecordingAppsService{
			Page: sdk.AppPage{
				Pagination: sdk.Pagination{
					TotalResults: 20,
					TotalPages:   10,
					PerPage:      2,
					Page:         0,
				},
				Apps: []sdk.App{
					{Id: "guid-a", Name: "name-a"},
					{Id: "guid-b", Name: "name-b"},
				},
			}}, expectedApps: &fakes.AppsRecorder{
			PerPage: 2,
			Page:    0,
		}},
		{desc: "lists apps perPage missing", method: "GET", path: "/api/apps?page=5", apps: &fakes.RecordingAppsService{}, expectedApps: &fakes.AppsRecorder{}},
		{desc: "lists apps perPage malformed", method: "GET", path: "/api/apps?perPage=a&page=5", apps: &fakes.RecordingAppsService{}, expectedApps: &fakes.AppsRecorder{}},
		{desc: "lists apps page malformed", method: "GET", path: "/api/apps?perPage=5&page=a", apps: &fakes.RecordingAppsService{}, expectedApps: &fakes.AppsRecorder{}},
		{desc: "error while listing apps", method: "GET", path: "/api/apps?perPage=5&page=2", apps: &fakes.RecordingAppsService{
			Err: someErr,
		}, expectedApps: &fakes.AppsRecorder{
			PerPage: 5,
			Page:    2,
		}},
		{desc: "gets pipeline", method: "GET", path: "/api/pipelines/guid-a", pipelines: &fakes.RecordingPipelinesService{
			Pipeline: sdk.Pipeline{
				Id: "guid-a", Name: "name-a",
			},
		}, expectedPipelines: &fakes.PipelinesRecorder{PipelineId: "guid-a"}},
		{desc: "error while getting pipeline", method: "GET", path: "/api/pipelines/guid-a", pipelines: &fakes.RecordingPipelinesService{
			Err: someErr,
		}, expectedPipelines: &fakes.PipelinesRecorder{
			PipelineId: "guid-a",
		}},
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
		{desc: "lists pipelines perPage missing", method: "GET", path: "/api/pipelines?page=5", pipelines: &fakes.RecordingPipelinesService{}, expectedPipelines: &fakes.PipelinesRecorder{}},
		{desc: "lists pipelines perPage empty", method: "GET", path: "/api/pipelines?perPage=", pipelines: &fakes.RecordingPipelinesService{}, expectedPipelines: &fakes.PipelinesRecorder{}},
		{desc: "lists pipelines perPage malformed", method: "GET", path: "/api/pipelines?perPage=a&page=5", pipelines: &fakes.RecordingPipelinesService{}, expectedPipelines: &fakes.PipelinesRecorder{}},
		{desc: "lists pipelines page malformed", method: "GET", path: "/api/pipelines?perPage=5&page=a", pipelines: &fakes.RecordingPipelinesService{}, expectedPipelines: &fakes.PipelinesRecorder{}},
		{desc: "error while listing pipelines", method: "GET", path: "/api/pipelines?perPage=5&page=2", pipelines: &fakes.RecordingPipelinesService{
			Err: someErr,
		}, expectedPipelines: &fakes.PipelinesRecorder{
			PerPage: 5,
			Page:    2,
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
		{desc: "error while getting pipeline runs", method: "GET", path: "/api/pipelines/pipeline-a/runs", pipelines: &fakes.RecordingPipelinesService{
			Err: someErr,
		}, expectedPipelines: &fakes.PipelinesRecorder{
			PipelineId: "pipeline-a",
			ToExcl:     someTime,
			Limit:      10,
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
		{
			desc:   "gets team",
			method: "GET",
			path:   "/api/teams/team-a",
			teams: &fakes.RecordingTeamsService{
				Team: teams.Team{
					Id: "team-id",
					TeamSettings: teams.TeamSettings{
						Name:        "team-name",
						Description: "team-desc",
						Access: teams.AccessGroups{
							Admin:  []string{"a"},
							Member: []string{"b"},
							Viewer: []string{"c"},
						},
					},
				},
			},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "deletes team",
			method: "DELETE",
			path:   "/api/teams/team-a",
			teams:  &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "updates team",
			method: "PUT",
			path:   "/api/teams/team-a",
			body: `
{
	"name": "test-name"
}
`,
			teams: &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId:   "team-a",
				TeamData: teams.TeamSettings{Name: "test-name"},
			},
		},
		{
			desc:   "creates team",
			method: "POST",
			path:   "/api/teams/team-a",
			body: `
{
	"name": "test-name"
}
`,
			teams: &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId:   "team-a",
				TeamData: teams.TeamSettings{Name: "test-name"},
			},
		},
		{
			desc:   "lists teams",
			method: "GET",
			path:   "/api/teams?perPage=5&page=2",
			teams: &fakes.RecordingTeamsService{
				Page: teams.TeamPage{Pagination: sdk.Pagination{
					TotalResults: 124,
					TotalPages:   5214,
					PerPage:      123,
					Page:         521,
				}, Teams: []teams.Team{{
					Id: "team-id",
					TeamSettings: teams.TeamSettings{
						Name:        "team-name",
						Description: "team-desc",
						Access: teams.AccessGroups{
							Admin:  []string{"a"},
							Member: []string{"b"},
							Viewer: []string{"c"},
						},
					},
				}}},
			},
			expectedTeams: &fakes.TeamsRecorder{
				PerPage: 5,
				Page:    2,
			},
		},
		{
			desc:          "listing teams: perPage missing",
			method:        "GET",
			path:          "/api/teams?page=2",
			teams:         &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{},
		},
		{
			desc:          "listing teams: perPage malformed",
			method:        "GET",
			path:          "/api/teams?perPage=abc&page=2",
			teams:         &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{},
		},
		{
			desc:   "listing teams: page missing",
			method: "GET",
			path:   "/api/teams?perPage=2",
			teams: &fakes.RecordingTeamsService{
				Page: teams.TeamPage{Pagination: sdk.Pagination{
					TotalResults: 124,
					TotalPages:   5214,
					PerPage:      123,
					Page:         521,
				}, Teams: []teams.Team{{
					Id: "team-id",
					TeamSettings: teams.TeamSettings{
						Name:        "team-name",
						Description: "team-desc",
						Access: teams.AccessGroups{
							Admin:  []string{"a"},
							Member: []string{"b"},
							Viewer: []string{"c"},
						},
					},
				}}},
			},
			expectedTeams: &fakes.TeamsRecorder{
				PerPage: 2,
				Page:    0,
			},
		},
		{
			desc:          "listing teams: page malformed",
			method:        "GET",
			path:          "/api/teams?perPage=5&page=abc",
			teams:         &fakes.RecordingTeamsService{},
			expectedTeams: &fakes.TeamsRecorder{},
		},
		{
			desc:   "listing teams: internal error",
			method: "GET",
			path:   "/api/teams?perPage=5&page=2",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{
				PerPage: 5,
				Page:    2,
			},
		},
		{
			desc:   "get team: internal error",
			method: "GET",
			path:   "/api/teams/team-a",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "update team: internal error",
			method: "PUT",
			path:   "/api/teams/team-a",
			body:   "{}",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "update team: malformed input",
			method: "PUT",
			path:   "/api/teams/team-a",
			body:   "abc",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{},
		},
		{
			desc:   "create team: internal error",
			method: "POST",
			path:   "/api/teams/team-a",
			body:   "{}",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "create team: malformed input",
			method: "POST",
			path:   "/api/teams/team-a",
			body:   "abc",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{},
		},
		{
			desc:   "delete team: internal error",
			method: "DELETE",
			path:   "/api/teams/team-a",
			teams: &fakes.RecordingTeamsService{
				Err: someErr,
			},
			expectedTeams: &fakes.TeamsRecorder{
				TeamId: "team-a",
			},
		},
		{
			desc:   "list groups",
			method: "GET",
			path:   "/api/groups",
			groups: &fakeGroups.RecordingGroupsService{
				ByProvider: map[string]groups.ProviderWithGroups{
					"provider-a": {
						Provider: "provider",
						Name:     "provider-name",
						Groups: []sdk.Group{
							{
								Id:   "group-a",
								Name: "group",
							},
						},
					},
				},
			},
			expectedGroups: &fakeGroups.GroupsRecorder{},
		},
		{
			desc:   "list groups error",
			method: "GET",
			path:   "/api/groups",
			groups: &fakeGroups.RecordingGroupsService{
				Err: someErr,
			},
			expectedGroups: &fakeGroups.GroupsRecorder{},
		},
		{
			desc:   "start websocket app",
			method: "GET",
			headers: http.Header{
				"Connection":            []string{"upgrade"},
				"Upgrade":               []string{"websocket"},
				"Sec-Websocket-Version": []string{"13"},
				"Sec-Websocket-Key":     []string{"abc"},
			},
			path: "/api/apps/app-a/live",
		},
		{
			desc:   "start websocket app with missing headers",
			method: "GET",
			path:   "/api/apps/app-a/live",
		},
	}

	currentTime = func() time.Time {
		return someTime
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			h := New(service.Core{
				Apps:      test.apps,
				Pipelines: test.pipelines,
				Teams:     test.teams,
				Groups:    test.groups,
			}, &fakes.PipeViz{}, Opts{
				DevConfig: config.DevConfig{DisableAuth: true},
			})

			testHttp(tt, h, test.method, test.path, test.body, test.headers)

			if test.expectedPipelines != nil && !cmp.Equal(*test.expectedPipelines, test.pipelines.Record) {
				tt.Errorf("pipeline records don't match:%s\n", cmp.Diff(*test.expectedPipelines, test.pipelines.Record))
			}

			if test.expectedApps != nil && !cmp.Equal(*test.expectedApps, test.apps.Record) {
				tt.Errorf("app records don't match:%s\n", cmp.Diff(*test.expectedApps, test.apps.Record))
			}

			if test.expectedTeams != nil && !cmp.Equal(*test.expectedTeams, test.teams.Record) {
				tt.Errorf("team records don't match:%s\n", cmp.Diff(*test.expectedTeams, test.teams.Record))
			}

			if test.expectedGroups != nil && !cmp.Equal(*test.expectedGroups, test.groups.Record) {
				tt.Errorf("group records don't match:%s\n", cmp.Diff(*test.expectedGroups, test.groups.Record))
			}
		})
	}

}

func TestDisableWebsocketXSRF(t *testing.T) {
	f := disableWebsocketXSRF(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))
	r, _ := http.NewRequest("GET", "/", nil)
	f.ServeHTTP(nil, r)

	r.Header.Set("Upgrade", "websocket")

	resp := &testResponseWriter{header: map[string][]string{}}
	f.ServeHTTP(resp, r)

	if resp.code != 403 {
		t.Error("expected 403")
	}
	expected := "{\"status\":403,\"error\":\"XSRF-TOKEN cookie not set\"}"
	respString := resp.String()
	if respString != expected {
		t.Errorf("mismatch: %s\n", cmp.Diff(expected, respString))
	}

	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: "secret"})
	resp = &testResponseWriter{header: map[string][]string{}}
	f.ServeHTTP(resp, r)

	if r.Header.Get("X-XSRF-TOKEN") != "secret" {
		t.Error("expected X-XSRF-TOKEN header from cookie")
	}
}

func testHttp(tt *testing.T, h http.Handler, method string, path string, body string, headers http.Header) {
	s := httptest.NewServer(h)
	defer s.Close()

	w := &hijackingResponseRecorder{httptest.NewRecorder()}
	r := httptest.NewRequest(method, s.URL+path, bytes.NewBuffer([]byte(body)))
	if headers != nil {
		r.Header = headers
	}

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

type testResponseWriter struct {
	bytes.Buffer
	header http.Header
	code   int
}

func (t *testResponseWriter) Header() http.Header {
	return t.header
}

func (t *testResponseWriter) WriteHeader(statusCode int) {
	t.code = statusCode
}

type hijackingResponseRecorder struct {
	*httptest.ResponseRecorder
}

func (h *hijackingResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := bufio.NewReadWriter(bufio.NewReader(&bytes.Buffer{}), bufio.NewWriter(&bytes.Buffer{}))
	return &fakeConn{}, rw, nil
}

type fakeConn struct{}

func (f *fakeConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (f *fakeConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (f *fakeConn) Close() error {
	return nil
}

func (f *fakeConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (f *fakeConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (f *fakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (f *fakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (f *fakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
