package api

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

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
