package pipeviz

import (
	"bytes"
	"io"
)

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	Id    int
	Label string
	Class string
}

type Edge struct {
	From  int
	To    int
	Class string
}

type PipeViz interface {
	Write(graph Graph, w io.Writer)
	Generate(graph Graph) []byte
}

func New() PipeViz {
	return &pipeViz{}
}

type pipeViz struct{}

func (p *pipeViz) Write(graph Graph, w io.Writer) {
	ir := generateLayout(graph)
	generateSvg(ir, w)
}

func (p *pipeViz) Generate(graph Graph) []byte {
	var buf bytes.Buffer
	p.Write(graph, &buf)
	return buf.Bytes()
}

func Generate(graph Graph, w io.Writer) {
	p := pipeViz{}
	p.Write(graph, w)
}
