package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

func (a *api) listAppsPaginated(w http.ResponseWriter, r *http.Request) {
	perPage, err := mustQueryInt(r, "perPage")
	if err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	page, err := defaultQueryInt(r, "page", 0)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}

	apps, err := a.db.ListAppsPaginated(perPage, page)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, apps)
}

func (a *api) getApp(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	app, err := a.db.GetApp(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, app)
}