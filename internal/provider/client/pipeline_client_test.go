package client

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http/httptest"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

func TestGetPipeline(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		state       fakePipelineProvider
		expectedErr error
	}{
		{desc: "returns pipeline", id: "id-a", state: fakePipelineProvider{
			pipeline: sdk.Pipeline{Id: "id-a", Name: "name-a"}},
		},
		{desc: "returns not found err", id: "not-exist", state: fakePipelineProvider{
			err: sdk.ErrNotFound,
		}, expectedErr: sdk.ErrNotFound},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewPipelineProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewPipelineProviderClient(s.URL, nil)

			pipeline, err := c.GetPipeline(test.id)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if test.state.recordedId != test.id {
				tt.Errorf("\nwanted id: %s\ngot: %s", test.id, test.state.recordedId)
			}

			if !cmp.Equal(test.state.pipeline, pipeline) {
				tt.Errorf("\ndiff between pipelines\n%s\n", cmp.Diff(test.state.pipeline, pipeline))
			}
		})
	}

}

func TestListPipelines(t *testing.T) {
	tests := []struct {
		desc        string
		state       fakePipelineProvider
		expectedErr error
	}{
		{desc: "returns pipelines", state: fakePipelineProvider{
			pipelines: []sdk.Pipeline{
				{Id: "id-a", Name: "name-a"},
				{Id: "id-b", Name: "name-b"},
			},
		},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewPipelineProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewPipelineProviderClient(s.URL, nil)

			pipelines, err := c.ListPipelines()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.state.pipelines, pipelines) {
				tt.Errorf("\ndiff between pipelines\n%s\n", cmp.Diff(test.state.pipelines, pipelines))
			}
		})
	}
}

func TestGetHistory(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		before      time.Time
		limit       int
		state       fakePipelineProvider
		expectedErr error
	}{
		{desc: "returns history", id: "a", before: someTime, limit: 20, state: fakePipelineProvider{
			history: sdk.PipelineStatusList{
				{
					Started: someTime,
					Steps: []sdk.StepRun{
						{
							StepId:  0,
							Status:  sdk.StatusSuccess,
							Started: someTime,
							Ended:   someTime,
						},
					},
				},
			},
		},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewPipelineProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewPipelineProviderClient(s.URL, nil)

			pipelines, err := c.GetHistory(test.id, test.before, test.limit)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if test.state.recordedId != test.id {
				tt.Errorf("\nwanted id: %s\ngot: %s", test.id, test.state.recordedId)
			}
			if test.state.recordedBefore != test.before {
				tt.Errorf("\nwanted before: %s\ngot: %s", test.before, test.state.recordedBefore)
			}
			if test.state.recordedLimit != test.limit {
				tt.Errorf("\nwanted limit: %d\ngot: %d", test.limit, test.state.recordedLimit)
			}

			if !cmp.Equal(test.state.history, pipelines) {
				tt.Errorf("\ndiff between pipelines\n%s\n", cmp.Diff(test.state.history, pipelines))
			}
		})
	}
}

func TestListUpdates(t *testing.T) {
	tests := []struct {
		desc            string
		since           time.Time
		state           fakePipelineProvider
		expectedUpdates sdk.PipelineUpdates
		expectedErr     error
	}{
		{desc: "returns updates", since: someTime.Add(-time.Minute), state: fakePipelineProvider{
			history: []sdk.PipelineStatus{
				{
					Started: someTime,
					Steps: []sdk.StepRun{
						{
							StepId: 0,
							Status: sdk.StatusSuccess,
						},
					},
				},
				{
					Steps: []sdk.StepRun{
						{
							StepId:  0,
							Started: someTime,
							Status:  sdk.StatusSuccess,
						},
					},
				},
				{
					Steps: []sdk.StepRun{
						{
							StepId: 0,
							Ended:  someTime,
							Status: sdk.StatusSuccess,
						},
					},
				},
				{
					Started: someTime.Add(-time.Hour),
				},
			},
			versions: sdk.PipelineVersionList{
				{
					Created: someTime,
				},
				{
					Created: someTime.Add(-time.Hour),
				},
			},
		},
			expectedUpdates: sdk.PipelineUpdates{
				Runs: sdk.PipelineStatusList{
					{
						Started: someTime,
						Steps: []sdk.StepRun{
							{
								StepId: 0,
								Status: sdk.StatusSuccess,
							},
						},
					},
					{
						Steps: []sdk.StepRun{
							{
								StepId:  0,
								Started: someTime,
								Status:  sdk.StatusSuccess,
							},
						},
					},
					{
						Steps: []sdk.StepRun{
							{
								StepId: 0,
								Ended:  someTime,
								Status: sdk.StatusSuccess,
							},
						},
					},
				},
				Versions: sdk.PipelineVersionList{
					{
						Created: someTime,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			handler := sdk.NewPipelineProviderHandler(&test.state)
			s := httptest.NewServer(handler)
			defer s.Close()

			c := NewPipelineProviderClient(s.URL, nil)

			updates, err := c.ListUpdates(test.since)
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err: %v\ngot: %v", test.expectedErr, err)
			}

			if !cmp.Equal(test.expectedUpdates, updates) {
				tt.Errorf("\ndiff between pipelines\n%s\n", cmp.Diff(test.expectedUpdates, updates))
			}
		})
	}
}

type fakePipelineProvider struct {
	err            error
	pipelines      []sdk.Pipeline
	pipeline       sdk.Pipeline
	recordedId     string
	recordedBefore time.Time
	recordedLimit  int
	history        sdk.PipelineStatusList
	versions       sdk.PipelineVersionList
}

func (f *fakePipelineProvider) ListUpdates(since time.Time) (sdk.PipelineUpdates, error) {
	updates := sdk.PipelineUpdates{}

	for _, version := range f.versions {
		if version.Created.After(since) {
			updates.Versions = append(updates.Versions, version)
		}
	}

	for _, status := range f.history {
		if status.Started.After(since) {
			updates.Runs = append(updates.Runs, status)
			continue
		}

		for _, step := range status.Steps {
			if step.Started.After(since) || step.Ended.After(since) {
				updates.Runs = append(updates.Runs, status)
				break
			}
		}
	}

	return updates, nil
}

func (f *fakePipelineProvider) ListPipelines() ([]sdk.Pipeline, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.pipelines, nil
}

func (f *fakePipelineProvider) GetPipeline(id string) (sdk.Pipeline, error) {
	f.recordedId = id
	if f.err != nil {
		return sdk.Pipeline{}, f.err
	}

	return f.pipeline, nil
}

func (f *fakePipelineProvider) GetHistory(id string, before time.Time, limit int) (sdk.PipelineStatusList, error) {
	f.recordedId = id
	f.recordedBefore = before
	f.recordedLimit = limit
	if f.err != nil {
		return nil, f.err
	}

	return f.history, nil
}
