package sdk

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

var instances = map[string]AppInstances{
	"a": {
		{
			State: "stopped",
			Since: someTime,
		},
	}}

func TestAppInstances(t *testing.T) {
	tests := []struct {
		desc           string
		state          map[string]AppInstances
		err            error
		method         string
		path           string
		expectedStatus int
		expectedResp   response
	}{
		{desc: "returns app instances", state: instances, method: "GET", path: "/instances/a", expectedStatus: http.StatusOK, expectedResp: response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"since": "2006-01-01T15:00:00Z", "state": "stopped"},
			},
		}},
		{desc: "returns 404 for non-existent", state: instances, method: "GET", path: "/instances/dont-exist", expectedStatus: http.StatusNotFound, expectedResp: response{
			Status: http.StatusNotFound,
			Err:    "not found",
		}},
		{desc: "returns 5xx for other errors", state: instances, err: ErrInternal, method: "GET", path: "/instances/a", expectedStatus: http.StatusInternalServerError, expectedResp: response{
			Status: http.StatusInternalServerError,
			Err:    "internal error occurred",
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := httptest.NewRecorder()
			p := &fakeInstancesProvider{state: test.state}
			p.err = test.err
			handler := NewAppInstancesProviderHandler(p)
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
		})
	}

}

type fakeInstancesProvider struct {
	state map[string]AppInstances
	err   error
}

func (f *fakeInstancesProvider) GetAppInstances(id string) (AppInstances, error) {
	if f.err != nil {
		return nil, f.err
	}

	if f.state[id] == nil {
		return nil, ErrNotFound
	}
	return f.state[id], nil
}
