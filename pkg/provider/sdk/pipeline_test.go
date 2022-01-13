package sdk

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestVersionListSortLen(t *testing.T) {
	l := PipelineVersionList{PipelineVersion{}, PipelineVersion{}}

	if l.Len() != 2 {
		t.Error("len is not 2")
	}
}

func TestVersionListSortLess(t *testing.T) {
	tests := []struct {
		desc     string
		list     PipelineVersionList
		i, j     int
		expected bool
	}{
		{desc: "is less", list: PipelineVersionList{
			PipelineVersion{Created: someTime},
			PipelineVersion{Created: someTime.Add(1 * time.Minute)},
		}, i: 0, j: 1, expected: true},
		{desc: "is not less", list: PipelineVersionList{
			PipelineVersion{Created: someTime},
			PipelineVersion{Created: someTime.Add(1 * time.Minute)},
		}, i: 1, j: 0, expected: false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			if test.list.Less(test.i, test.j) != test.expected {
				tt.Error("less result is not as expected")
			}
		})
	}
}

func TestVersionListSortSwap(t *testing.T) {
	tests := []struct {
		desc  string
		list  PipelineVersionList
		after PipelineVersionList
		i, j  int
	}{
		{desc: "swaps", list: PipelineVersionList{
			PipelineVersion{PipelineId: "a"},
			PipelineVersion{PipelineId: "b"},
		}, i: 0, j: 1, after: PipelineVersionList{
			PipelineVersion{PipelineId: "b"},
			PipelineVersion{PipelineId: "a"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			test.list.Swap(test.i, test.j)
			if !cmp.Equal(test.list, test.after) {
				tt.Errorf("diff: %s\n", cmp.Diff(test.after, test.list))
			}
		})
	}
}

func TestVersionAt(t *testing.T) {
	list := PipelineVersionList{
		PipelineVersion{Created: someTime.Add(-10 * time.Minute), PipelineId: "c"},
		PipelineVersion{Created: someTime.Add(-5 * time.Minute), PipelineId: "b"},
		PipelineVersion{Created: someTime, PipelineId: "a"},
	}

	tests := []struct {
		desc     string
		at       time.Time
		expected string
	}{
		{"gets previous for exact match", someTime, "b"},
		{"gets previous for in between", someTime.Add(-6 * time.Minute), "c"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			version := list.VersionAt(test.at)
			if version.PipelineId != test.expected {
				tt.Errorf("expected '%s', was '%s'", test.expected, version.PipelineId)
			}
		})
	}
}

func TestStatusListSortLen(t *testing.T) {
	l := PipelineStatusList{PipelineStatus{}, PipelineStatus{}}

	if l.Len() != 2 {
		t.Error("len is not 2")
	}
}

func TestStatusListSortLess(t *testing.T) {
	tests := []struct {
		desc     string
		list     PipelineStatusList
		i, j     int
		expected bool
	}{
		{desc: "is less", list: PipelineStatusList{
			PipelineStatus{Started: someTime},
			PipelineStatus{Started: someTime.Add(1 * time.Minute)},
		}, i: 0, j: 1, expected: true},
		{desc: "is not less", list: PipelineStatusList{
			PipelineStatus{Started: someTime},
			PipelineStatus{Started: someTime.Add(1 * time.Minute)},
		}, i: 1, j: 0, expected: false},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			if test.list.Less(test.i, test.j) != test.expected {
				tt.Error("less result is not as expected")
			}
		})
	}
}

func TestStatusListSortSwap(t *testing.T) {
	tests := []struct {
		desc  string
		list  PipelineStatusList
		after PipelineStatusList
		i, j  int
	}{
		{desc: "swaps", list: PipelineStatusList{
			PipelineStatus{PipelineId: "a"},
			PipelineStatus{PipelineId: "b"},
		}, i: 0, j: 1, after: PipelineStatusList{
			PipelineStatus{PipelineId: "b"},
			PipelineStatus{PipelineId: "a"},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			test.list.Swap(test.i, test.j)
			if !cmp.Equal(test.list, test.after) {
				tt.Errorf("diff: %s\n", cmp.Diff(test.after, test.list))
			}
		})
	}
}

func TestStatusListFold(t *testing.T) {
	tests := []struct {
		desc     string
		list     PipelineStatusList
		expected PipelineStatus
	}{
		{
			desc:     "folds starting time",
			list:     PipelineStatusList{PipelineStatus{Started: someTime}, PipelineStatus{Started: someTime.Add(10 * time.Minute)}},
			expected: PipelineStatus{Started: someTime.Add(10 * time.Minute), Steps: []StepRun{}},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			status := test.list.Fold()
			if !cmp.Equal(status, test.expected) {
				tt.Errorf("diff: %s\n", cmp.Diff(test.expected, status))
			}
		})
	}
}
