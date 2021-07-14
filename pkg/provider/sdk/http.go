package sdk

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type ProviderConfig struct {
	Apps      AppProvider
	Pipelines PipelineProvider
}

func ListenAndServe(addr string, p ProviderConfig) error {
	h := mux.NewRouter()

	if p.Pipelines != nil {
		h.PathPrefix("/apps").Handler(NewAppProviderHandler(p.Apps))
	}

	if p.Pipelines != nil {
		h.PathPrefix("/pipelines").Handler(NewPipelineProviderHandler(p.Pipelines))
	}

	return http.ListenAndServe(addr, h)
}

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

func defaultQueryTime(r *http.Request, queryKey string, defaultValue time.Time) (time.Time, error) {
	valueStr := r.FormValue(queryKey)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := time.Parse(time.RFC3339, valueStr)
	if err != nil {
		return time.Time{}, err
	}

	return value, nil
}
