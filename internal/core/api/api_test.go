package api

import (
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

func TestName(t *testing.T) {
	tests := []struct {
		desc   string
		state  sdk.AppPage
		method string
		path   string
	}{
		{"lists apps", sdk.AppPage{
			TotalResults: 20,
			TotalPages:   10,
			PerPage:      2,
			Page:         5,
			Apps: []sdk.App{
				{Id: "guid-a", Name: "name-a"},
				{Id: "guid-b", Name: "name-b"},
			},
		}, "GET", "/apps?perPage=2&page=5"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			h := New(&fakeDb{test.state})

			s := httptest.NewServer(h)
			defer s.Close()

			w := httptest.NewRecorder()

			r := httptest.NewRequest(test.method, s.URL+test.path, nil)
			h.ServeHTTP(w, r)

			res, _ := httputil.DumpResponse(w.Result(), true)
			approvals.UseFolder("testdata")
			approvals.UseReporter(reporters.NewGoLandReporter())
			approvals.VerifyString(tt, string(res))
		})
	}

}

type fakeDb struct {
	page sdk.AppPage
}

func (f *fakeDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	if f.page.PerPage == perPage && f.page.Page == page {
		return f.page, nil
	}
	return sdk.AppPage{}, sdk.ErrPageExceeded
}
