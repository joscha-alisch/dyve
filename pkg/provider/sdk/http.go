package sdk

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)
import "github.com/gorilla/mux"

func ListenAndServeAppProvider(addr string, p AppProvider) error {
	return http.ListenAndServe(addr, NewAppProviderHandler(p))
}

func NewAppProviderHandler(p AppProvider) http.Handler {
	h := &appProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/apps", h.listApps)
	h.HandleFunc("/apps/{id:[0-9a-z-]+}", h.getApp)

	return h
}

type appProviderHandler struct {
	*mux.Router

	p AppProvider
}

type response struct {
	Status int         `json:"status"`
	Err    string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

func (h *appProviderHandler) listApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.p.ListApps()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, apps)
}

func (h *appProviderHandler) getApp(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetApp(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, app)
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
