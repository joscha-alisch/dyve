package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func New(db database.Database) http.Handler {
	a := &api{
		Router: mux.NewRouter(),
		db:     db,
	}

	a.Path("/api/apps").Queries("perPage", "").HandlerFunc(a.listAppsPaginated)
	a.Path("/api/apps/{id:[0-9a-z-]+}").HandlerFunc(a.getApp)
	a.Path("/api/pipelines").Queries("perPage", "").HandlerFunc(a.listPipelinesPaginated)
	a.Path("/api/pipelines/{id:[0-9a-z-]+}").HandlerFunc(a.getPipeline)

	return a
}

type api struct {
	*mux.Router
	db database.Database
}

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

func (a *api) listPipelinesPaginated(w http.ResponseWriter, r *http.Request) {
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

	pipelines, err := a.db.ListPipelinesPaginated(perPage, page)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, pipelines)
}

func (a *api) getPipeline(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pipeline, err := a.db.GetPipeline(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, pipeline)
}

func mustQueryInt(r *http.Request, queryKey string) (int, error) {
	valueStr := r.FormValue(queryKey)
	if valueStr == "" {
		return 0, errExpectedQueryParamMissing
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func defaultQueryInt(r *http.Request, queryKey string, defaultValue int) (int, error) {
	valueStr := r.FormValue(queryKey)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, nil
}

var errExpectedQueryParamMissing = errors.New("query parameter was expected but is missing")

type response struct {
	Status int         `json:"status"`
	Err    string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

func respondOk(w http.ResponseWriter, result interface{}) {
	respond(w, response{
		Status: http.StatusOK,
		Result: result,
	})
}

func respondErr(w http.ResponseWriter, code int, err error) {
	respond(w, response{
		Status: code,
		Err:    err.Error(),
	})
}

func respond(w http.ResponseWriter, r response) {
	b, err := json.Marshal(r)
	if err != nil {
		log.Error().Interface("response", r).Err(err).Msg("error marshalling response")
	}

	w.WriteHeader(r.Status)
	_, err = w.Write(b)
	if err != nil {
		log.Error().Interface("response", r).Err(err).Msg("error writing response")
	}
}
