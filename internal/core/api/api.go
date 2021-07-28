package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"net/http"
)

func New(db database.Database, pipeGen pipeviz.PipeViz) http.Handler {
	a := &api{
		Router:  mux.NewRouter(),
		db:      db,
		pipeGen: pipeGen,
	}

	a.Path("/api/apps").Queries("perPage", "").HandlerFunc(a.listAppsPaginated)
	a.Path("/api/apps/{id:[0-9a-z-]+}").HandlerFunc(a.getApp)

	a.Path("/api/pipelines").Queries("perPage", "").HandlerFunc(a.listPipelinesPaginated)
	a.Path("/api/pipelines/{id:[0-9a-z-]+}/status").HandlerFunc(a.getPipelineStatus)
	a.Path("/api/pipelines/{id:[0-9a-z-]+}/runs").HandlerFunc(a.listPipelineRuns)
	a.Path("/api/pipelines/{id:[0-9a-z-]+}").HandlerFunc(a.getPipeline)

	return a
}

type api struct {
	*mux.Router
	db      database.Database
	pipeGen pipeviz.PipeViz
}
