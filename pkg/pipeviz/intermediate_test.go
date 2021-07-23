package pipeviz

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestIR(t *testing.T) {
	tests := []struct {
		desc     string
		g        Graph
		expected layouted
	}{
		{"one node", Graph{Nodes: []Node{
			{0, "build", "step"},
		}}, layouted{
			Rows: 1, Columns: 1,
			Boxes: []box{
				{0, 0, 0, 0, Node{0, "build", "step"}},
			},
		}},
		{"two parallel nodes", Graph{Nodes: []Node{
			{0, "build", "step"},
			{1, "test", "step"},
		}}, layouted{
			Rows: 2, Columns: 1,
			Boxes: []box{
				{0, 0, 0, 0, Node{0, "build", "step"}},
				{0, 0, 1, 1, Node{1, "test", "step"}},
			},
		}},
		{"two sequential nodes", Graph{Nodes: []Node{
			{0, "build", "step"},
			{1, "test", "step"},
		}, Edges: []Edge{
			{0, 1, "input"},
		}}, layouted{
			Rows: 1, Columns: 2,
			Boxes: []box{
				{0, 0, 0, 0, Node{0, "build", "step"}},
				{1, 1, 0, 0, Node{1, "test", "step"}},
			},
		}},
		{"combined after parallel", Graph{Nodes: []Node{
			{0, "build", "step"},
			{1, "test", "step"},
			{2, "deploy", "step"},
		}, Edges: []Edge{
			{0, 2, "input"},
			{1, 2, "input"},
		}}, layouted{
			Rows: 2, Columns: 2,
			Boxes: []box{
				{0, 0, 0, 0, Node{0, "build", "step"}},
				{0, 0, 1, 1, Node{1, "test", "step"}},
				{1, 1, 0, 1, Node{2, "deploy", "step"}},
			},
		}},
		{"over multiple columns", Graph{Nodes: []Node{
			{0, "build", "step"},
			{1, "test", "step"},
			{3, "build2", "step"},
			{2, "deploy", "step"},
		}, Edges: []Edge{
			{0, 1, "input"},
			{1, 2, "input"},
			{3, 2, "input"},
		}}, layouted{
			Rows: 2, Columns: 3,
			Boxes: []box{
				{0, 0, 1, 1, Node{0, "build", "step"}},
				{1, 1, 1, 1, Node{1, "test", "step"}},
				{2, 2, 0, 1, Node{2, "deploy", "step"}},
				{0, 1, 0, 0, Node{3, "build2", "step"}},
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			res := generateLayout(test.g)
			if !cmp.Equal(test.expected, res, cmp.AllowUnexported(layouted{})) {
				tt.Errorf("found diff:\n%s", cmp.Diff(test.expected, res))
			}
		})
	}
}
