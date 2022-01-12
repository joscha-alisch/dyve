package sdk

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

var routes = map[string]AppRoutes{
	"a": {
		{
			Host:    "host",
			Path:    "path",
			AppPort: 900,
		},
	}}

func TestAppRouting(t *testing.T) {
	tests := []struct {
		desc           string
		state          map[string]AppRoutes
		method         string
		path           string
		expectedStatus int
		expectedResp   response
	}{
		{"returns app routing", routes, "GET", "/routing/a", http.StatusOK, response{
			Status: http.StatusOK,
			Result: map[string]interface{}{"routes": []interface{}{
				map[string]interface{}{"host": "host", "path": "path", "appPort": float64(900)},
			}},
		}},
		{"returns 404 for non-existent", routes, "GET", "/routing/dont-exist", http.StatusNotFound, response{
			Status: http.StatusNotFound,
			Err:    "not found",
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := httptest.NewRecorder()
			handler := NewAppRoutingProviderHandler(&fakeRoutingProvider{state: test.state})
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

type fakeRoutingProvider struct {
	state map[string]AppRoutes
}

func (f *fakeRoutingProvider) GetAppRouting(id string) (AppRouting, error) {
	for appId, routes := range f.state {
		if appId == id {
			return AppRouting{Routes: routes}, nil
		}
	}
	return AppRouting{}, ErrNotFound
}
