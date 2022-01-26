package pipelines

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

var somePipeline = sdk.Pipeline{
	Id:   "pipeline-a",
	Name: "pipeline",
}
var somePipelineStatus = sdk.PipelineStatus{
	PipelineId: "pipeline-a",
}
var somePipelineVersion = sdk.PipelineVersion{
	PipelineId: "pipeline-a",
	Created:    someTime.Add(-1 * time.Minute),
}
var somePagination = sdk.Pagination{
	TotalResults: 123,
	TotalPages:   123,
	PerPage:      123,
	Page:         123,
}
var someErr = errors.New("some error")
var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

type DecodableFunc func(target interface{}) error

func (f DecodableFunc) Decode(dec interface{}) error {
	return f(dec)
}

func TestService_ListPipelinesPaginated(t *testing.T) {
	tests := []struct {
		desc        string
		perPage     int
		page        int
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.PipelinePage
		expectedErr error
	}{
		{
			desc:    "gets instances",
			perPage: 5,
			page:    2,
			db: &db.Database{
				ReturnPagination: func(pagination *sdk.Pagination) {
					*pagination = somePagination
				},
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*sdk.Pipeline) = somePipeline
						return nil
					}))
				},
			},
			expected: sdk.PipelinePage{
				Pagination: somePagination,
				Pipelines:  []sdk.Pipeline{somePipeline},
			},
			recorded: []db.DatabaseRecord{{
				Collection: "pipelines",
				PerPage:    5,
				Page:       2,
			}},
		},
		{
			desc:    "error while getting group",
			perPage: 5,
			page:    2,
			db: &db.Database{
				Err: someErr,
			},
			expected:    sdk.PipelinePage{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListPipelinesPaginated(test.perPage, test.page)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_GetPipeline(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.Pipeline
		expectedErr error
	}{
		{
			desc: "gets instances",
			id:   "pipeline-a",
			db: &db.Database{
				Return: func(target interface{}) {
					*target.(*sdk.Pipeline) = somePipeline
				},
			},
			expected: somePipeline,
			recorded: []db.DatabaseRecord{{
				Collection: "pipelines",
				Id:         "pipeline-a",
			}},
		},
		{
			desc: "error while getting group",
			id:   "pipeline-a",
			db: &db.Database{
				Err: someErr,
			},
			expected:    sdk.Pipeline{},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.GetPipeline(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_ListPipelineRuns(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		fromIncl    time.Time
		toExcl      time.Time
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.PipelineStatusList
		expectedErr error
	}{
		{
			desc:     "lists pipeline runs",
			id:       "pipeline-a",
			fromIncl: someTime,
			toExcl:   someTime.Add(2 * time.Minute),
			db: &db.Database{
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*sdk.PipelineStatus) = somePipelineStatus
						return nil
					}))
				},
			},
			expected: sdk.PipelineStatusList{
				somePipelineStatus,
			},
			recorded: []db.DatabaseRecord{{
				Collection: "pipeline_runs",
				Filter: bson.M{
					"pipelineId": bson.M{
						"$eq": "pipeline-a",
					},
					"started": bson.M{
						"$lt":  someTime.Add(2 * time.Minute),
						"$gte": someTime,
					},
				},
			}},
		},
		{
			desc:     "error while getting group",
			id:       "pipeline-a",
			fromIncl: someTime,
			toExcl:   someTime.Add(2 * time.Minute),
			db: &db.Database{
				Err: someErr,
			},
			expected:    nil,
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListPipelineRuns(test.id, test.fromIncl, test.toExcl)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_ListPipelineRunsLimit(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		toExcl      time.Time
		limit       int
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.PipelineStatusList
		expectedErr error
	}{
		{
			desc:   "lists pipeline runs",
			id:     "pipeline-a",
			toExcl: someTime.Add(2 * time.Minute),
			limit:  20,
			db: &db.Database{
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*sdk.PipelineStatus) = somePipelineStatus
						return nil
					}))
				},
			},
			expected: sdk.PipelineStatusList{
				somePipelineStatus,
			},
			recorded: []db.DatabaseRecord{{
				Collection: "pipeline_runs",
				Limit:      20,
				Sort: bson.M{
					"started": -1,
				},
				Filter: bson.M{
					"pipelineId": bson.M{
						"$eq": "pipeline-a",
					},
					"started": bson.M{
						"$lt": someTime.Add(2 * time.Minute),
					},
				},
			}},
		},
		{
			desc:   "error while getting group",
			id:     "pipeline-a",
			toExcl: someTime.Add(2 * time.Minute),
			limit:  20,
			db: &db.Database{
				Err: someErr,
			},
			expected:    nil,
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListPipelineRunsLimit(test.id, test.toExcl, test.limit)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_ListPipelineVersions(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		fromIncl    time.Time
		toExcl      time.Time
		db          *db.Database
		recorded    []db.DatabaseRecord
		expected    sdk.PipelineVersionList
		expectedErr error
	}{
		{
			desc:     "lists pipeline versions",
			id:       "pipeline-a",
			fromIncl: someTime,
			toExcl:   someTime.Add(2 * time.Minute),
			db: &db.Database{
				Return: func(target interface{}) {
					*target.(*sdk.PipelineVersion) = somePipelineVersion
				},
				ReturnEach: func(each func(decodable database.Decodable) error) {
					_ = each(DecodableFunc(func(target interface{}) error {
						*target.(*sdk.PipelineVersion) = somePipelineVersion
						return nil
					}))
				},
			},
			expected: sdk.PipelineVersionList{
				somePipelineVersion,
				somePipelineVersion,
			},
			recorded: []db.DatabaseRecord{{
				Collection: "pipeline_versions",
				Filter: bson.M{
					"pipelineId": "pipeline-a",
					"created": bson.M{
						"$lt":  someTime.Add(2 * time.Minute),
						"$gte": someTime,
					},
				},
			}, {
				Collection: "pipeline_versions",
				Sort: bson.M{
					"created": -1,
				},
				Filter: bson.M{
					"pipelineId": bson.M{
						"$eq": "pipeline-a",
					},
					"created": bson.M{
						"$lt": someTime,
					},
				},
			}},
		},
		{
			desc:     "error while listing versions",
			id:       "pipeline-a",
			fromIncl: someTime,
			toExcl:   someTime.Add(2 * time.Minute),
			db: &db.Database{
				Err: someErr,
			},
			expected:    nil,
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			res, err := s.ListPipelineVersions(test.id, test.fromIncl, test.toExcl)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.expected, res) {
				tt.Errorf("results mismatch: %s\n", cmp.Diff(test.expected, res))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}

		})
	}
}

func TestService_AddPipelineRuns(t *testing.T) {
	tests := []struct {
		desc        string
		provider    string
		runs        sdk.PipelineStatusList
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:     "adds pipeline runs",
			provider: "provider-a",
			runs:     []sdk.PipelineStatus{somePipelineStatus},
			db:       &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection: "pipeline_runs",
				Updates: map[string]interface{}{
					"pipeline-a": somePipelineStatus,
				},
				Filters: map[string]interface{}{
					"pipeline-a": bson.M{
						"pipelineId": "pipeline-a",
						"provider":   "provider-a",
						"started":    time.Time{},
					},
				},
			}},
		},
		{
			desc:     "error while updating pipelines",
			provider: "provider-a",
			runs:     []sdk.PipelineStatus{somePipelineStatus},
			db: &db.Database{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.AddPipelineRuns(test.provider, test.runs)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}

func TestService_AddPipelineVersions(t *testing.T) {
	tests := []struct {
		desc        string
		provider    string
		versions    sdk.PipelineVersionList
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:     "adds pipeline versions",
			provider: "provider-a",
			versions: sdk.PipelineVersionList{somePipelineVersion},
			db:       &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection: "pipeline_versions",
				Updates: map[string]interface{}{
					"pipeline-a2006-01-01T14:59:00Z": somePipelineVersion,
				},
				Filters: map[string]interface{}{
					"pipeline-a2006-01-01T14:59:00Z": bson.M{
						"pipelineId": "pipeline-a",
						"provider":   "provider-a",
						"created":    someTime.Add(-1 * time.Minute),
					},
				},
			}},
		},
		{
			desc:     "error while updating pipelines",
			provider: "provider-a",
			versions: sdk.PipelineVersionList{somePipelineVersion},
			db: &db.Database{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.AddPipelineVersions(test.provider, test.versions)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}
func TestService_UpdatePipelines(t *testing.T) {
	tests := []struct {
		desc        string
		provider    string
		pipelines   []sdk.Pipeline
		db          *db.Database
		recorded    []db.DatabaseRecord
		expectedErr error
	}{
		{
			desc:      "updates pipelines",
			provider:  "provider-a",
			pipelines: []sdk.Pipeline{somePipeline},
			db:        &db.Database{},
			recorded: []db.DatabaseRecord{{
				Collection: "pipelines",
				Provider:   "provider-a",
				Updates: map[string]interface{}{
					"pipeline-a": somePipeline,
				},
			}},
		},
		{
			desc:      "error while updating pipelines",
			provider:  "provider-a",
			pipelines: []sdk.Pipeline{somePipeline},
			db: &db.Database{
				Err: someErr,
			},
			expectedErr: someErr,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			recorder := &db.DatabaseRecorder{}
			test.db.Recorder = recorder
			s := NewService(test.db)
			err := s.UpdatePipelines(test.provider, test.pipelines)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("errors mismatch: %s\n", cmp.Diff(test.expectedErr, err))
			}

			if !cmp.Equal(test.recorded, recorder.Records) {
				tt.Errorf("recorded mismatch: %s\n", cmp.Diff(test.recorded, recorder.Records))
			}
		})
	}
}
