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
		method         string
		path           string
		expectedStatus int
		expectedResp   response
	}{
		{"returns app instances", instances, "GET", "/instances/a", http.StatusOK, response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"since": "2006-01-01T15:00:00Z", "state": "stopped"},
			},
		}},
		{"returns 404 for non-existent", instances, "GET", "/instances/dont-exist", http.StatusNotFound, response{
			Status: http.StatusNotFound,
			Err:    "not found",
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := httptest.NewRecorder()
			handler := NewAppInstancesProviderHandler(&fakeInstancesProvider{state: test.state})
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
}

func (f *fakeInstancesProvider) GetAppInstances(id string) (AppInstances, error) {
	if f.state[id] == nil {
		return nil, ErrNotFound
	}
	return f.state[id], nil
}
