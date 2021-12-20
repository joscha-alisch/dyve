package fakes

import (
	"bytes"
	"encoding/json"
	"github.com/joscha-alisch/dyve/pkg/pipeviz"
	"io"
)

type PipeViz struct {
}

func (f *PipeViz) Write(graph pipeviz.Graph, w io.Writer) {
	b, _ := json.Marshal(graph)
	_, _ = w.Write(b)
}

func (f *PipeViz) Generate(graph pipeviz.Graph) []byte {
	var buf bytes.Buffer
	buf.Write([]byte("fake svg: "))
	f.Write(graph, &buf)
	return buf.Bytes()
}
