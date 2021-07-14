package sdk

import (
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")
var otherTime, _ = time.Parse(time.RFC3339, "2006-01-01T13:00:00Z")

var pipelines = []Pipeline{
	{Id: "a", Name: "name-a", Definition: PipelineDefinition{
		Connections: []PipelineConnection{
			{From: 0, To: 1, Manual: true},
		},
		Steps: []PipelineStep{
			{Name: "test", Id: 0},
			{Name: "build", Id: 1, AppDeployments: []string{"app-a"}},
		},
	}},
	{Id: "840e560f-38d3-460e-be23-8677a4539f35", Name: "name-b"},
}

var runs = []PipelineRun{
	{
		PipelineId: "a",
		Start:      someTime,
		Steps: []StepRun{
			{
				StepId:  0,
				Started: someTime,
				Ended:   someTime.Add(time.Minute),
				Status:  StatusSuccess,
			},
		},
	},
}

func TestPipelines(t *testing.T) {
	currentTime = func() time.Time {
		return someTime
	}

	tests := []struct {
		desc                  string
		state                 fakePipelineProvider
		method                string
		path                  string
		expectedStatus        int
		expectedResp          response
		expectedRecordedId    string
		expectedRecordedTime  time.Time
		expectedRecordedLimit int
	}{
		{desc: "returns pipelines", state: fakePipelineProvider{
			pipelines: pipelines,
		}, method: "GET", path: "/pipelines", expectedStatus: http.StatusOK, expectedResp: response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"id": "a", "name": "name-a", "definition": map[string]interface{}{
					"connections": []interface{}{
						map[string]interface{}{"from": float64(0), "to": float64(1), "manual": true},
					},
					"steps": []interface{}{
						map[string]interface{}{"name": "test", "id": float64(0)},
						map[string]interface{}{"name": "build", "id": float64(1), "appDeployments": []interface{}{"app-a"}},
					},
				}},
				map[string]interface{}{"id": "840e560f-38d3-460e-be23-8677a4539f35", "name": "name-b", "definition": map[string]interface{}{}},
			},
		}},
		{desc: "returns pipelines with trailing slash", state: fakePipelineProvider{
			pipelines: []Pipeline{},
		}, method: "GET", path: "/pipelines/", expectedStatus: http.StatusOK, expectedResp: response{
			Status: http.StatusOK,
			Result: []interface{}{},
		}},
		{desc: "returns pipelines error", state: fakePipelineProvider{
			err: errors.New("some error that should not be broadcast"),
		}, method: "GET", path: "/pipelines/", expectedStatus: http.StatusInternalServerError, expectedResp: response{
			Status: http.StatusInternalServerError,
			Err:    ErrInternal.Error(),
		}},
		{desc: "returns pipeline", state: fakePipelineProvider{
			pipeline: pipelines[0],
		}, method: "GET", path: "/pipelines/a", expectedRecordedId: "a", expectedStatus: http.StatusOK, expectedResp: response{
			Status: http.StatusOK,
			Result: map[string]interface{}{"id": "a", "name": "name-a", "definition": map[string]interface{}{
				"connections": []interface{}{
					map[string]interface{}{"from": float64(0), "to": float64(1), "manual": true},
				},
				"steps": []interface{}{
					map[string]interface{}{"name": "test", "id": float64(0)},
					map[string]interface{}{"name": "build", "id": float64(1), "appDeployments": []interface{}{"app-a"}},
				},
			}},
		}},
		{desc: "returns pipeline not found", state: fakePipelineProvider{
			err: ErrNotFound,
		}, method: "GET", path: "/pipelines/not-exist", expectedStatus: http.StatusNotFound, expectedResp: response{
			Status: http.StatusNotFound,
			Err:    ErrNotFound.Error(),
		}},
		{desc: "returns pipeline internal error", state: fakePipelineProvider{
			err: errors.New("error that should not be returned"),
		}, method: "GET", path: "/pipelines/not-exist", expectedStatus: http.StatusInternalServerError, expectedResp: response{
			Status: http.StatusInternalServerError,
			Err:    ErrInternal.Error(),
		}},
		{desc: "returns history with defaults", state: fakePipelineProvider{
			history: runs,
		}, method: "GET", path: "/pipelines/a/history",
			expectedRecordedId:    "a",
			expectedRecordedTime:  someTime,
			expectedRecordedLimit: 10,
			expectedStatus:        http.StatusOK, expectedResp: response{
				Status: http.StatusOK,
				Result: []interface{}{
					map[string]interface{}{
						"pipelineId": "a", "start": someTime.Format(time.RFC3339), "steps": []interface{}{
							map[string]interface{}{
								"stepId":  float64(0),
								"started": someTime.Format(time.RFC3339),
								"ended":   someTime.Add(time.Minute).Format(time.RFC3339),
								"status":  StatusSuccess,
							},
						},
					},
				},
			}},
		{desc: "returns history with query params", state: fakePipelineProvider{
			history: []PipelineRun{},
		}, method: "GET", path: "/pipelines/a/history?before=2006-01-01T13:00:00Z&limit=20",
			expectedRecordedId:    "a",
			expectedRecordedTime:  otherTime,
			expectedRecordedLimit: 20,
			expectedStatus:        http.StatusOK, expectedResp: response{
				Status: http.StatusOK,
				Result: []interface{}{},
			}},
		{desc: "returns history pipeline not found", state: fakePipelineProvider{
			err: ErrNotFound,
		}, method: "GET", path: "/pipelines/not-exist", expectedStatus: http.StatusNotFound, expectedResp: response{
			Status: http.StatusNotFound,
			Err:    ErrNotFound.Error(),
		}},
		{desc: "returns history pipeline internal error", state: fakePipelineProvider{
			err: errors.New("error that should not be returned"),
		}, method: "GET", path: "/pipelines/a/history", expectedStatus: http.StatusInternalServerError, expectedResp: response{
			Status: http.StatusInternalServerError,
			Err:    ErrInternal.Error(),
		}},
		{desc: "returns history pipeline limit malformed error", state: fakePipelineProvider{},
			method: "GET", path: "/pipelines/a/history?limit=abc", expectedStatus: http.StatusBadRequest, expectedResp: response{
				Status: http.StatusBadRequest,
				Err:    ErrQueryLimitMalformed.Error(),
			}},
		{desc: "returns history pipeline since malformed error", state: fakePipelineProvider{
			err: ErrQuerySinceMalformed,
		}, method: "GET", path: "/pipelines/a/history?before=abc", expectedStatus: http.StatusBadRequest, expectedResp: response{
			Status: http.StatusBadRequest,
			Err:    ErrQuerySinceMalformed.Error(),
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := httptest.NewRecorder()
			handler := NewPipelineProviderHandler(&test.state)
			handler.ServeHTTP(r, httptest.NewRequest(test.method, test.path, nil))
			res := r.Result()
			if res.StatusCode != test.expectedStatus {
				tt.Errorf("\nwanted stats %v\n   got %v", test.expectedStatus, res.StatusCode)
			}

			resp := response{}
			_ = json.NewDecoder(res.Body).Decode(&resp)
			if !cmp.Equal(test.expectedResp, resp) {
				tt.Errorf("\ndiff between responses: \n%s\n", cmp.Diff(test.expectedResp, resp))
			}

			if test.state.recordedId != test.expectedRecordedId {
				tt.Errorf("recorded id: %s, expected: %s\n", test.state.recordedId, test.expectedRecordedId)
			}

			if test.state.recordedTime != test.expectedRecordedTime {
				tt.Errorf("recorded time: %s, expected: %s\n", test.state.recordedTime, test.expectedRecordedTime)
			}

			if test.state.recordedLimit != test.expectedRecordedLimit {
				tt.Errorf("recorded id: %d, expected: %d\n", test.state.recordedLimit, test.expectedRecordedLimit)
			}
		})
	}

}

type fakePipelineProvider struct {
	pipelines []Pipeline
	pipeline  Pipeline
	history   []PipelineRun
	err       error

	recordedId    string
	recordedTime  time.Time
	recordedLimit int
}

func (f *fakePipelineProvider) ListPipelines() ([]Pipeline, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.pipelines, nil
}

func (f *fakePipelineProvider) GetPipeline(id string) (Pipeline, error) {
	if f.err != nil {
		return Pipeline{}, f.err
	}

	f.recordedId = id
	return f.pipeline, nil
}

func (f *fakePipelineProvider) GetHistory(id string, since time.Time, limit int) ([]PipelineRun, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.recordedId = id
	f.recordedTime = since
	f.recordedLimit = limit

	return f.history, nil
}
