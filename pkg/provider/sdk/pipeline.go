package sdk

import (
	"time"
)

type PipelinePage struct {
	Pagination
	Pipelines []Pipeline
}

type Pipeline struct {
	Id         string             `json:"id"`
	Name       string             `json:"name"`
	Definition PipelineDefinition `json:"definition"`
}

type PipelineDefinition struct {
	Steps       []PipelineStep       `json:"steps,omitempty"`
	Connections []PipelineConnection `json:"connections,omitempty"`
}

type PipelineStep struct {
	Name           string   `json:"name"`
	Id             int      `json:"id"`
	AppDeployments []string `json:"appDeployments,omitempty"`
}

type PipelineConnection struct {
	From   int  `json:"from"`
	To     int  `json:"to"`
	Manual bool `json:"manual"`
}

type PipelineRun struct {
	PipelineId string    `json:"pipelineId"`
	Start      time.Time `json:"start"`
	Steps      []StepRun `json:"steps"`
}

type StepRun struct {
	StepId  int            `json:"stepId"`
	Status  PipelineStatus `json:"status"`
	Started time.Time      `json:"started"`
	Ended   time.Time      `json:"ended"`
}

type PipelineStatus string

const (
	StatusSuccess = "succeeded"
	StatusFailure = "failed"
	StatusRunning = "running"
	StatusAborted = "aborted"
	StatusPending = "pending"
)

type PipelineProvider interface {
	ListPipelines() ([]Pipeline, error)
	GetPipeline(id string) (Pipeline, error)
	GetHistory(id string, before time.Time, limit int) ([]PipelineRun, error)
}
