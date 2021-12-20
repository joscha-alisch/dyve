package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
	"sort"
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

	pipelines, err := a.core.Pipelines.ListPipelinesPaginated(perPage, page)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, pipelines)
}

func (a *api) getPipeline(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pipeline, err := a.core.Pipelines.GetPipeline(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	respondOk(w, pipeline)
}

func (a *api) listPipelineRuns(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	before, _ := defaultQueryTime(r, "before", currentTime())
	limit, _ := defaultQueryInt(r, "limit", 10)

	runs, err := a.core.Pipelines.ListPipelineRunsLimit(id, before, limit)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	if len(runs) == 0 {
		respondOk(w, []pipelineStatus{})
		return
	}

	sort.Sort(runs)

	var res []pipelineStatus
	versions, err := a.core.Pipelines.ListPipelineVersions(id, runs[0].Started, before)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	for _, run := range runs {
		version := versions.VersionAt(run.Started)
		if version.PipelineId == "" {
			continue
		}

		svg := a.toSvg(version, run)

		res = append(res, pipelineStatus{
			PipelineStatus: run,
			Svg:            string(svg),
		})
	}

	respondOk(w, res)
}

func (a *api) getPipelineStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	pipeline, err := a.core.Pipelines.GetPipeline(id)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	runs, err := a.core.Pipelines.ListPipelineRuns(id, pipeline.Current.Created, currentTime())
	if err != nil {
		respondErr(w, http.StatusInternalServerError, sdk.ErrInternal)
		return
	}

	status := runs.Fold()
	svg := a.toSvg(pipeline.Current, status)

	respondOk(w, pipelineStatus{
		PipelineStatus: status,
		Svg:            string(svg),
	})
}

func (a *api) toSvg(version sdk.PipelineVersion, status sdk.PipelineStatus) []byte {
	g := pipeviz.Graph{
		Edges: make([]pipeviz.Edge, len(version.Definition.Connections)),
	}
	nodes := make(map[int]*pipeviz.Node)

	for _, step := range version.Definition.Steps {
		nodes[step.Id] = &pipeviz.Node{
			Id:    step.Id,
			Label: step.Name,
			Class: "",
		}
	}
	for i, connection := range version.Definition.Connections {
		g.Edges[i] = pipeviz.Edge{From: connection.From, To: connection.To}
	}

	for _, step := range status.Steps {
		nodes[step.StepId].Class = string(step.Status)
	}

	g.Nodes = make([]pipeviz.Node, len(nodes))
	i := 0
	for _, node := range nodes {
		g.Nodes[i] = *node
		i++
	}

	sort.Slice(g.Nodes, func(i, j int) bool {
		return g.Nodes[i].Id < g.Nodes[j].Id
	})

	return a.pipeGen.Generate(g)
}
