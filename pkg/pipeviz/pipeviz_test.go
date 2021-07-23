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
				{Id: 0, Label: "build", Class: "step"},
				{Id: 1, Label: "test", Class: "step"},
			},
		}},
		{"two inputs", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "step"},
				{Id: 1, Label: "test-a", Class: "step"},
				{Id: 2, Label: "deploy", Class: "step"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 1, To: 2},
			},
		}},
		{"two outputs", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "step"},
				{Id: 1, Label: "notify", Class: "step"},
				{Id: 2, Label: "deploy", Class: "step"},
			},
			Edges: []Edge{
				{From: 0, To: 1},
				{From: 0, To: 2},
			},
		}},
		{"multiple columns", Graph{
			Nodes: []Node{
				{Id: 0, Label: "build-a", Class: "step"},
				{Id: 1, Label: "build-b", Class: "step"},
				{Id: 2, Label: "test", Class: "step"},
				{Id: 3, Label: "deploy", Class: "step"},
			},
			Edges: []Edge{
				{From: 0, To: 2},
				{From: 2, To: 3},
				{From: 1, To: 3},
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

	Create(g, &buf)

	return buf.Bytes()
}
