package fakes

import (
	"github.com/joscha-alisch/dyve/internal/core/pipelines"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

type RecordingPipelinesService struct {
	Err       error
	Pipeline  sdk.Pipeline
	Page      sdk.PipelinePage
	Pipelines []sdk.Pipeline
	Runs      sdk.PipelineStatusList
	Versions  sdk.PipelineVersionList
	Record    PipelinesRecorder
}

type PipelinesRecorder struct {
	PerPage    int
	Page       int
	PipelineId string
	FromIncl   time.Time
	ToExcl     time.Time
	Limit      int
	ProviderId string
	Pipelines  []sdk.Pipeline
	Runs       sdk.PipelineStatusList
	Versions   sdk.PipelineVersionList
}

func (s *RecordingPipelinesService) EnsureIndices() error {
	return nil
}

func (s *RecordingPipelinesService) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	s.Record.PerPage = perPage
	s.Record.Page = page

	if s.Err != nil {
		return sdk.PipelinePage{}, s.Err
	}
	return s.Page, nil
}

func (s *RecordingPipelinesService) GetPipeline(id string) (sdk.Pipeline, error) {
	s.Record.PipelineId = id
	if s.Err != nil {
		return sdk.Pipeline{}, s.Err
	}
	return s.Pipeline, nil
}

func (s *RecordingPipelinesService) ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error) {
	s.Record.PipelineId = id
	s.Record.FromIncl = fromIncl
	s.Record.ToExcl = toExcl

	if s.Err != nil {
		return nil, s.Err
	}
	return s.Runs, nil
}

func (s *RecordingPipelinesService) ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error) {
	s.Record.PipelineId = id
	s.Record.ToExcl = toExcl
	s.Record.Limit = limit

	if s.Err != nil {
		return nil, s.Err
	}
	return s.Runs, nil
}

func (s *RecordingPipelinesService) ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error) {
	s.Record.PipelineId = id
	s.Record.FromIncl = fromIncl
	s.Record.ToExcl = toExcl

	if s.Err != nil {
		return nil, s.Err
	}
	return s.Versions, nil
}

func (s *RecordingPipelinesService) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	s.Record.ProviderId = providerId
	s.Record.Pipelines = pipelines

	if s.Err != nil {
		return s.Err
	}
	return nil
}

func (s *RecordingPipelinesService) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	s.Record.ProviderId = providerId
	s.Record.Runs = runs

	if s.Err != nil {
		return s.Err
	}
	return nil
}

func (s *RecordingPipelinesService) AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error {
	s.Record.ProviderId = providerId
	s.Record.Versions = versions

	if s.Err != nil {
		return s.Err
	}
	return nil
}

func NewPipelineMapping(p []pipelines.Pipeline, v []sdk.PipelineVersion, r []sdk.PipelineStatus) *MappingPipelinesService {
	m := &MappingPipelinesService{
		Pipelines: make(map[string]pipelines.Pipeline),
		Runs:      make(map[string]sdk.PipelineStatusList),
		Versions:  make(map[string]sdk.PipelineVersionList),
	}
	for _, pipeline := range p {
		m.Pipelines[pipeline.Id] = pipeline
	}

	for _, version := range v {
		m.Versions[version.PipelineId] = append(m.Versions[version.PipelineId], version)
	}
	for _, run := range r {
		m.Runs[run.PipelineId] = append(m.Runs[run.PipelineId], run)
	}

	return m
}

type MappingPipelinesService struct {
	Pipelines map[string]pipelines.Pipeline
	Runs      map[string]sdk.PipelineStatusList
	Versions  map[string]sdk.PipelineVersionList
}

func (m *MappingPipelinesService) EnsureIndices() error {
	panic("implement me")
}

func (m *MappingPipelinesService) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	panic("implement me")
}

func (m *MappingPipelinesService) GetPipeline(id string) (sdk.Pipeline, error) {
	return m.Pipelines[id].Pipeline, nil
}

func (m *MappingPipelinesService) ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error) {
	var res sdk.PipelineStatusList
	for _, status := range m.Runs[id] {
		if status.Started.After(fromIncl.Add(-1*time.Second)) && status.Started.Before(toExcl) {
			res = append(res, status)
		}
	}
	return res, nil
}

func (m *MappingPipelinesService) ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error) {
	var res sdk.PipelineStatusList
	p := m.Runs[id]
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Started.Before(toExcl) {
			res = append(res, p[i])
		}

		if len(res) >= limit {
			break
		}
	}
	return res, nil
}

func (m *MappingPipelinesService) ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MappingPipelinesService) UpdatePipelines(providerId string, pipelineList []sdk.Pipeline) error {
	for id, pipeline := range m.Pipelines {
		if pipeline.ProviderId == providerId {
			delete(m.Pipelines, id)
		}
	}

	for _, pipeline := range pipelineList {
		m.Pipelines[pipeline.Id] = pipelines.Pipeline{
			Pipeline:   pipeline,
			ProviderId: providerId,
		}
	}
	return nil
}

func (m *MappingPipelinesService) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	if m.Runs == nil {
		m.Runs = map[string]sdk.PipelineStatusList{}
	}

	for _, run := range runs {
		m.Runs[run.PipelineId] = append(m.Runs[run.PipelineId], run)
	}
	return nil
}

func (m *MappingPipelinesService) AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error {
	if m.Versions == nil {
		m.Versions = map[string]sdk.PipelineVersionList{}
	}

	for _, version := range versions {
		m.Versions[version.PipelineId] = append(m.Versions[version.PipelineId], version)
	}
	return nil
}
