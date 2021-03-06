package sdk

import (
	"time"
)

type PipelinePage struct {
	Pagination
	Pipelines []Pipeline `json:"pipelines" bson:"pipelines"`
}

type PipelineUpdates struct {
	Runs     PipelineStatusList  `json:"runs,omitempty" bson:"runs,omitempty"`
	Versions PipelineVersionList `json:"versions,omitempty" bson:"versions,omitempty"`
}
type Pipeline struct {
	Id      string          `json:"id" bson:"id"`
	Name    string          `json:"name" bson:"name"`
	Current PipelineVersion `json:"current" bson:"current"`
}

type PipelineVersionList []PipelineVersion

func (pl PipelineVersionList) Len() int {
	return len(pl)
}

func (pl PipelineVersionList) Less(i, j int) bool {
	return pl[i].Created.Before(pl[j].Created)
}

func (pl PipelineVersionList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl PipelineVersionList) VersionAt(t time.Time) PipelineVersion {
	for i := len(pl) - 1; i >= 0; i-- {
		if pl[i].Created.Before(t) {
			return pl[i]
		}
	}
	return PipelineVersion{}
}

type PipelineVersion struct {
	PipelineId string             `json:"pipelineId" bson:"pipelineId"`
	Created    time.Time          `json:"created" bson:"created"`
	Definition PipelineDefinition `json:"definition" bson:"definition"`
}

type PipelineDefinition struct {
	Steps       []PipelineStep       `json:"steps,omitempty" bson:"steps,omitempty"`
	Connections []PipelineConnection `json:"connections,omitempty" bson:"connections,omitempty"`
}

type PipelineStep struct {
	Name           string   `json:"name" bson:"name"`
	Id             int      `json:"id" bson:"id"`
	AppDeployments []string `json:"appDeployments,omitempty" bson:"appDeployments,omitempty"`
}

type PipelineConnection struct {
	From   int  `json:"from" bson:"from"`
	To     int  `json:"to" bson:"to"`
	Manual bool `json:"manual" bson:"manual"`
}

type PipelineStatusList []PipelineStatus

type PipelineStatus struct {
	PipelineId string    `json:"pipelineId" bson:"pipelineId"`
	Started    time.Time `json:"started" bson:"started"`
	Steps      []StepRun `json:"steps,omitempty" bson:"steps,omitempty"`
}

type StepRun struct {
	StepId  int        `json:"stepId" bson:"stepId"`
	Status  StepStatus `json:"status" bson:"status"`
	Started time.Time  `json:"started" bson:"started"`
	Ended   time.Time  `json:"ended" bson:"ended"`
}

type StepStatus string

const (
	StatusSuccess = "succeeded"
	StatusFailure = "failed"
	StatusRunning = "running"
	StatusAborted = "aborted"
	StatusPending = "pending"
)

type PipelineProvider interface {
	ListPipelines() ([]Pipeline, error)
	ListUpdates(since time.Time) (PipelineUpdates, error)
	GetPipeline(id string) (Pipeline, error)
	GetHistory(id string, before time.Time, limit int) (PipelineStatusList, error)
}

func (pl PipelineStatusList) Fold() PipelineStatus {
	if len(pl) == 1 {
		return pl[0]
	}

	s := PipelineStatus{}
	for _, status := range pl {
		s = s.Fold(status)
	}
	return s
}

func (pl PipelineStatusList) Len() int {
	return len(pl)
}

func (pl PipelineStatusList) Less(i, j int) bool {
	return pl[i].Started.Before(pl[j].Started)
}

func (pl PipelineStatusList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (p PipelineStatus) Fold(other PipelineStatus) PipelineStatus {
	p.PipelineId = other.PipelineId

	if p.Started.Before(other.Started) {
		p.Started = other.Started
	}

	m := make(map[int]StepRun, len(p.Steps))
	for _, step := range p.Steps {
		m[step.StepId] = step
	}

	for _, step := range other.Steps {
		if step.Started.After(m[step.StepId].Started) {
			m[step.StepId] = step
		}
	}

	p.Steps = make([]StepRun, len(m))
	i := 0
	for _, run := range m {
		p.Steps[i] = run
		i++
	}

	return p
}
