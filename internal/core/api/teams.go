package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

func (a *api) listTeamsPaginated(w http.ResponseWriter, r *http.Request) {
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

	teams, err := a.core.Teams.ListTeamsPaginated(perPage, page)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, teams)
}

func (a *api) getTeam(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	team, err := a.core.Teams.GetTeam(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, team)
}

func (a *api) upsertTeam(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := a.core.Teams.UpdateTeam(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, nil)
}

func (a *api) deleteTeam(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := a.core.Teams.DeleteTeam(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, nil)
}
