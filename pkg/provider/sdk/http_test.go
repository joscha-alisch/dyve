package sdk

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

var apps = []App{
	{Id: "a", Name: "name-a"},
	{Id: "b", Name: "name-b"},
	{Id: "c", Name: "name-c"},
	{Id: "d", Name: "name-d"},
	{Id: "e", Name: "name-e"},
	{Id: "f", Name: "name-f"},
	{Id: "g", Name: "name-g"},
	{Id: "h", Name: "name-h"},
	{Id: "i", Name: "name-i"},
	{Id: "j", Name: "name-j"},
	{Id: "840e560f-38d3-460e-be23-8677a4539f35", Name: "name-k"},
}

func TestName(t *testing.T) {
	tests := []struct {
		desc           string
		state          []App
		method         string
		path           string
		expectedStatus int
		expectedResp   response
	}{
		{"returns apps", apps, "GET", "/apps", http.StatusOK, response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"id": "a", "name": "name-a"},
				map[string]interface{}{"id": "b", "name": "name-b"},
				map[string]interface{}{"id": "c", "name": "name-c"},
				map[string]interface{}{"id": "d", "name": "name-d"},
				map[string]interface{}{"id": "e", "name": "name-e"},
				map[string]interface{}{"id": "f", "name": "name-f"},
				map[string]interface{}{"id": "g", "name": "name-g"},
				map[string]interface{}{"id": "h", "name": "name-h"},
				map[string]interface{}{"id": "i", "name": "name-i"},
				map[string]interface{}{"id": "j", "name": "name-j"},
				map[string]interface{}{"id": "840e560f-38d3-460e-be23-8677a4539f35", "name": "name-k"},
			},
		}},
		{"returns app", apps, "GET", "/apps/840e560f-38d3-460e-be23-8677a4539f35", http.StatusOK, response{
			Status: http.StatusOK,
			Result: map[string]interface{}{"id": "840e560f-38d3-460e-be23-8677a4539f35", "name": "name-k"},
		}},
		{"returns 404 for non-existent", apps, "GET", "/apps/dont-exist", http.StatusNotFound, response{
			Status: http.StatusNotFound,
			Err:    "not found",
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			r := httptest.NewRecorder()
			handler := NewAppProviderHandler(&fakeProvider{state: test.state})
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

type fakeProvider struct {
	state []App
}

func (f *fakeProvider) ListApps() ([]App, error) {
	return f.state, nil
}

func (f *fakeProvider) GetApp(id string) (App, error) {
	for _, app := range f.state {
		if app.Id == id {
			return app, nil
		}
	}
	return App{}, ErrNotFound
}
