package pipeviz

import (
	"bytes"
	"errors"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestApproval(t *testing.T) {
	tests := []struct {
		desc string
		g    Graph
	}{
		{"simple", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build", Class: "running"},
				{Id: 1, Label: "test", Class: "succeeded"},
			},
		}},
		{"two inputs", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "step"},
				{Id: 1, Label: "test-a", Class: "running"},
				{Id: 2, Label: "deploy", Class: "succeeded"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 1, To: 2},
			},
		}},
		{"two outputs", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "running"},
				{Id: 1, Label: "notify", Class: "succeeded"},
				{Id: 2, Label: "deploy", Class: "step"},
			},
			Edges: []Edge{
				{From: 0, To: 1},
				{From: 0, To: 2},
			},
		}},
		{"multiple columns", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "succeeded"},
				{Id: 1, Label: "build-b", Class: "running"},
				{Id: 2, Label: "test", Class: "succeeded"},
				{Id: 3, Label: "deploy", Class: "failed"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 2, To: 3},
				{From: 1, To: 3},
			},
		}},
		{"different col lengths", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "succeeded"},
				{Id: 1, Label: "test-a", Class: "running"},
				{Id: 2, Label: "deploy-a", Class: "running"},
				{Id: 3, Label: "build-b", Class: "failed"},
				{Id: 4, Label: "deploy-b", Class: "failed"},
				{Id: 5, Label: "notification", Class: "succeeded"},
			},
			Edges: []Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 2, To: 5},
				{From: 3, To: 4},
				{From: 4, To: 5},
			},
		}},
		{"cross link", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "failed"},
				{Id: 1, Label: "test-ab", Class: "running"},
				{Id: 2, Label: "deploy-a", Class: "succeeded"},
				{Id: 3, Label: "build-b", Class: "failed"},
				{Id: 4, Label: "deploy-b", Class: "succeeded"},
				{Id: 5, Label: "notification", Class: "failed"},
			},
			Edges: []Edge{
				{From: 0, To: 1},
				{From: 1, To: 2},
				{From: 2, To: 5},
				{From: 3, To: 4},
				{From: 4, To: 5},
				{From: 1, To: 4},
			},
		}},
		{"2 to 4", Graph{
			Nodes: []Node{
				{Id: 0, Label: "first-a", Class: "succeeded"},
				{Id: 1, Label: "first-b", Class: "running"},
				{Id: 2, Label: "second-a", Class: "succeeded"},
				{Id: 3, Label: "second-b", Class: "failed"},
				{Id: 4, Label: "second-c", Class: "succeeded"},
				{Id: 5, Label: "second-d", Class: "failed"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 0, To: 3},
				{From: 1, To: 4},
				{From: 1, To: 5},
			},
		}},
		{"2 to 5", Graph{
			Nodes: []Node{
				{Id: 0, Label: "first-a", Class: "failed"},
				{Id: 1, Label: "first-b", Class: "running"},
				{Id: 2, Label: "second-a", Class: "succeeded"},
				{Id: 3, Label: "second-b", Class: "failed"},
				{Id: 4, Label: "second-c", Class: "succeeded"},
				{Id: 5, Label: "second-d", Class: "succeeded"},
				{Id: 6, Label: "second-e", Class: "succeeded"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 0, To: 3},
				{From: 1, To: 4},
				{From: 1, To: 5},
				{From: 1, To: 6},
			},
		}},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			runApprovalTest(tt, test.desc, test.g)
		})
	}
}

func runApprovalTest(t *testing.T, testName string, g Graph) {
	fileName := strings.ReplaceAll(strings.ToLower(testName), " ", "_")
	acceptedName := fileName + ".accepted.svg"
	actualName := fileName + ".actual.svg"
	_, testFilePath, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFilePath)
	testDir = filepath.Join(testDir, "approval_tests")

	_ = os.Remove(filepath.Join(testDir, actualName))

	approved := loadApproved(filepath.Join(testDir, acceptedName))
	actual := createActual(g)

	if !cmp.Equal(approved, actual) {
		t.Errorf("actual does not equal approved: \n%s", cmp.Diff(string(approved), string(actual)))

		_ = ioutil.WriteFile(filepath.Join(testDir, actualName), actual, 0666)
	}
}

func loadApproved(fileName string) []byte {
	f, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return []byte{}
	}
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return b
}

func createActual(g Graph) []byte {
	var buf bytes.Buffer

	Generate(g, &buf)

	return buf.Bytes()
}
