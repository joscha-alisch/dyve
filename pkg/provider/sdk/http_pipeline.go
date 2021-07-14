package sdk

import (
	"errors"
	"net/http"
	"time"
)
import "github.com/gorilla/mux"

var currentTime = time.Now

func ListenAndServePipelineProvider(addr string, p PipelineProvider) error {
	return ListenAndServe(addr, ProviderConfig{Pipelines: p})
}

func NewPipelineProviderHandler(p PipelineProvider) http.Handler {
	h := &pipelineProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/pipelines", h.listPipelines)
	h.HandleFunc("/pipelines/", h.listPipelines)
	h.HandleFunc("/pipelines/{id:[0-9a-z-]+}", h.getPipeline)
	h.HandleFunc("/pipelines/{id:[0-9a-z-]+}/history", h.getHistory)

	return h
}

type pipelineProviderHandler struct {
	*mux.Router

	p PipelineProvider
}

func (h *pipelineProviderHandler) listPipelines(w http.ResponseWriter, r *http.Request) {
	apps, err := h.p.ListPipelines()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, ErrInternal)
		return
	}
	respondOk(w, apps)
}

func (h *pipelineProviderHandler) getPipeline(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetPipeline(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	respondOk(w, app)
}

func (h *pipelineProviderHandler) getHistory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	limit, err := defaultQueryInt(r, "limit", 10)
	if err != nil {
		respondErr(w, http.StatusBadRequest, ErrQueryLimitMalformed)
		return
	}

	since, err := defaultQueryTime(r, "before", currentTime())
	if err != nil {
		respondErr(w, http.StatusBadRequest, ErrQuerySinceMalformed)
		return
	}

	history, err := h.p.GetHistory(id, since, limit)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, ErrInternal)
		return
	}

	respondOk(w, history)
}
