package pipeviz

import (
	"gonum.org/v1/gonum/graph/simple"
)

type layouted struct {
	Rows    int
	Columns int
	Boxes   []box
}

type box struct {
	StartCol int
	EndCol   int
	StartRow int
	EndRow   int
	Node     Node
	Lines    []float64
}

func generateLayout(in Graph) layouted {
	nodes := make(map[int]Node)
	edges := make(map[int][]Edge)

	g := simple.NewDirectedGraph()
	for _, n := range in.Nodes {
		g.AddNode(simple.Node(n.Id))
		nodes[n.Id] = n
	}
	for _, edge := range in.Edges {
		g.SetEdge(g.NewEdge(simple.Node(edge.From), simple.Node(edge.To)))
		edges[edge.From] = append(edges[edge.From], edge)
	}

	layout := newGrid()

	for _, n := range in.Nodes {
		layout.addBox(n.Id)
	}
	for _, n := range in.Nodes {
		to := g.To(int64(n.Id))
		for to.Next() {
			layout.constrainRightOf(n.Id, int(to.Node().ID()))
		}
	}

	layout.layout()

	boxes := layout.getBoxes()
	var out []box

	for id, b := range boxes {
		b.Node = nodes[boxes[id].Node.Id]

		for _, edge := range edges[id] {
			start := float64(boxes[edge.To].StartRow)
			end := float64(boxes[edge.To].EndRow)

			target := start + (end-start)/2
			b.Lines = append(b.Lines, target)
		}

		out = append(out, b)
	}

	return layouted{
		Rows:    layout.height,
		Columns: layout.width,
		Boxes:   out,
	}
}
