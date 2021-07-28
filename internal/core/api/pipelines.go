package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
	"time"
)

var currentTime = time.Now

type pipelineStatus struct {
	sdk.PipelineStatus
	Svg string `json:"svg"`
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

func (a *api) listPipelineRuns(w http.ResponseWriter, r *http.Request) {

}

func (a *api) getPipelineStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pipeline, err := a.db.GetPipeline(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	runs, err := a.db.ListPipelineRuns(id, pipeline.Current.Created, currentTime())
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	g := pipeviz.Graph{
		Edges: make([]pipeviz.Edge, len(pipeline.Current.Definition.Connections)),
	}
	nodes := make(map[int]*pipeviz.Node)

	for _, step := range pipeline.Current.Definition.Steps {
		nodes[step.Id] = &pipeviz.Node{
			Id:    step.Id,
			Label: step.Name,
			Class: "",
		}
	}
	for i, connection := range pipeline.Current.Definition.Connections {
		g.Edges[i] = pipeviz.Edge{From: connection.From, To: connection.To}
	}

	status := runs.Fold()
	for _, step := range status.Steps {
		nodes[step.StepId].Class = string(step.Status)
	}

	g.Nodes = make([]pipeviz.Node, len(nodes))
	i := 0
	for _, node := range nodes {
		g.Nodes[i] = *node
		i++
	}

	svg := a.pipeGen.Generate(g)

	respondOk(w, pipelineStatus{
		PipelineStatus: status,
		Svg:            string(svg),
	})
}
