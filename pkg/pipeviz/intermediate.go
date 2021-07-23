package pipeviz

import (
	"gonum.org/v1/gonum/graph/simple"
	"sort"
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
}

type subGraph struct {
	nodes  map[int]bool
	inputs []int
}

func generateLayout(in Graph) layouted {
	nodes := make(map[int]Node)

	g := simple.NewDirectedGraph()
	for _, n := range in.Nodes {
		g.AddNode(simple.Node(n.Id))
		nodes[n.Id] = n
	}
	for _, edge := range in.Edges {
		g.SetEdge(g.NewEdge(simple.Node(edge.From), simple.Node(edge.To)))
	}

	root := make(map[int]bool)
	leaf := make(map[int]bool)

	for _, n := range in.Nodes {
		if g.To(int64(n.Id)).Len() == 0 {
			root[n.Id] = true
		}
		if g.From(int64(n.Id)).Len() == 0 {
			leaf[n.Id] = true
		}
	}

	var current = leaf
	var previous map[int]bool

	columns := make([]map[int]bool, 0)
	col := 0

	for len(current) > 0 {
		previous = make(map[int]bool)
		for n, _ := range current {
			if len(columns) <= col {
				columns = append(columns, make(map[int]bool, 0))
			}

			columns[col][n] = true

			prev := g.To(int64(n))
			for prev.Next() {
				previous[int(prev.Node().ID())] = true
			}
		}
		current = previous
		col++
	}

	for i := 0; i < len(columns)-1; i++ {
		for n, _ := range columns[i] {
			to := g.To(int64(n))
			previousColumnHasPrev := false
			for to.Next() {
				if columns[i+1][int(to.Node().ID())] {
					previousColumnHasPrev = true
					break
				}
			}

			if !previousColumnHasPrev {
				columns[i+1][n] = true
			}
		}
	}

	var colLengths []int
	for _, column := range columns {
		colLengths = append(colLengths, len(column))
	}
	gridRows := leastCommonMultiple(colLengths...)

	boxes := make(map[int]*box)

	currentCol := len(columns) - 1

	for col := 0; col < len(columns); col++ {
		column := columns[col]
		i := 0
		for node := range column {
			if boxes[node] == nil {
				boxes[node] = &box{
					StartCol: currentCol,
					EndCol:   currentCol,
					StartRow: 0,
					EndRow:   0,
					Node:     nodes[node],
				}
			} else {
				boxes[node].StartCol = currentCol
			}
			i++
		}

		currentCol--
	}

	var boxList []box

	for _, b := range boxes {
		boxList = append(boxList, *b)
	}

	sort.Slice(boxList, func(i, j int) bool {
		return boxList[i].Node.Id < boxList[j].Node.Id
	})

	return layouted{
		Rows:    gridRows,
		Columns: len(columns),
		Boxes:   boxList,
	}
}

// greatest common divisor (GCD) via Euclidean algorithm
func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func leastCommonMultiple(integers ...int) int {
	if len(integers) < 2 {
		return integers[0]
	}

	a := integers[0]
	b := integers[1]

	result := a * b / greatestCommonDivisor(a, b)

	for i := 2; i < len(integers); i++ {
		result = leastCommonMultiple(result, integers[i])
	}

	return result
}
