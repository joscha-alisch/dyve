package sdk

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var apps = []App{
	{"a", "name-a"},
	{"b", "name-b"},
	{"c", "name-c"},
	{"d", "name-d"},
	{"e", "name-e"},
	{"f", "name-f"},
	{"g", "name-g"},
	{"h", "name-h"},
	{"i", "name-i"},
	{"j", "name-j"},
	{"840e560f-38d3-460e-be23-8677a4539f35", "name-k"},
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
		{"returns search result", apps, "GET", "/search?term=name-a&limit=10", http.StatusOK, response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "a", "name": "name-a"}},
			},
		}},
		{"limits search results", apps, "GET", "/search?term=name&limit=2", http.StatusOK, response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "a", "name": "name-a"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "b", "name": "name-b"}},
			},
		}},
		{"errors on invalid limit", apps, "GET", "/search?term=name&limit=bla", http.StatusBadRequest, response{
			Status: http.StatusBadRequest,
			Err:    "query param 'limit' is not an integer",
		}},
		{"errors on missing search term", apps, "GET", "/search", http.StatusBadRequest, response{
			Status: http.StatusBadRequest,
			Err:    "query param 'term' is required",
		}},
		{"uses default limit of 10", apps, "GET", "/search?term=name", http.StatusOK, response{
			Status: http.StatusOK,
			Result: []interface{}{
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "a", "name": "name-a"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "b", "name": "name-b"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "c", "name": "name-c"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "d", "name": "name-d"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "e", "name": "name-e"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "f", "name": "name-f"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "g", "name": "name-g"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "h", "name": "name-h"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "i", "name": "name-i"}},
				map[string]interface{}{"score": 0.4, "app": map[string]interface{}{"id": "j", "name": "name-j"}},
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
		{"lists apps paged", apps, "GET", "/apps?perPage=2", http.StatusOK, response{
			Status: http.StatusOK,
			Result: map[string]interface{}{
				"totalPages": float64(6),
				"totalResults": float64(11),
				"perPage": float64(2),
				"page": float64(0),
				"apps": []interface{}{
					map[string]interface{}{"id": "a", "name": "name-a"},
					map[string]interface{}{"id": "b", "name": "name-b"},
				},
			},
		}},
		{"fails for malformed page", apps, "GET", "/apps?perPage=2&page=abc", http.StatusBadRequest, response{
			Status: http.StatusBadRequest,
			Err: ErrQueryPageMalformed.Error(),
		}},
		{"fails for malformed perPage", apps, "GET", "/apps?perPage=abc&page=1", http.StatusBadRequest, response{
			Status: http.StatusBadRequest,
			Err: ErrQueryPerPageMalformed.Error(),
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

func (f *fakeProvider) ListAppsPaged(perPage int, page int) (AppPage, error) {
	total := len(f.state)
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	cursor := perPage * page
	if cursor > total {
		return AppPage{}, ErrPageExceeded
	}

	until := math.Min(float64(total-1), float64(cursor+perPage))
	return AppPage{
		TotalPages: pages,
		TotalResults: total,
		PerPage: perPage,
		Page: page,
		Apps: f.state[cursor:int(until)],
	}, nil
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

func (f *fakeProvider) Search(term string, limit int) ([]AppSearchResult, error) {
	var matches []App
	for _, app := range f.state {
		if len(matches) >= limit {
			break
		}
		if strings.Contains(app.Name, term) {
			matches = append(matches, app)
		}
	}

	var res []AppSearchResult
	for _, match := range matches {
		res = append(res, AppSearchResult{
			App:   match,
			Score: 0.4,
		})
	}
	return res, nil
}
