package pipeviz

import (
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

func Create(graph Graph, w io.Writer) {
	ir := generateLayout(graph)
	generateSvg(ir, w)
}
