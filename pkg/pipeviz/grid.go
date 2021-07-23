package pipeviz

type grid struct {
	constraints      []constraint
	constraintSource map[int][]constraint
	constraintTarget map[int][]constraint

	width  int
	height int
	rows   map[int]map[int]bool
	cols   map[int]map[int]bool
	boxes  map[int]position
}

func newGrid() *grid {
	return &grid{
		constraints:      nil,
		constraintSource: make(map[int][]constraint),
		constraintTarget: make(map[int][]constraint),
		width:            0,
		height:           0,
		rows:             make(map[int]map[int]bool),
		cols:             make(map[int]map[int]bool),
		boxes:            nil,
	}
}

type constraint struct {
	t           constraintType
	src, target int
}

type constraintType int

const (
	constraintRightOf constraintType = iota
)

type position struct {
	row    int
	col    int
	width  int
	height int
}

func (g *grid) addBox(id int) {
	g.boxes[id] = position{}
}

func (g *grid) constrainRightOf(src, target int) {
	c := constraint{
		t:      constraintRightOf,
		src:    src,
		target: target,
	}
	g.constraints = append(g.constraints, c)

	g.constraintSource[src] = append(g.constraintSource[src], c)
	g.constraintTarget[target] = append(g.constraintTarget[target], c)
}

func (g *grid) layout() {

}

func (g *grid) getBoxes() []box {
	var b []box

	for id, p := range g.boxes {
		b = append(b, box{
			StartCol: p.col,
			EndCol:   p.col + (p.width - 1),
			StartRow: p.row,
			EndRow:   p.row + (p.height - 1),
			Node: Node{
				Id: id,
			},
		})
	}
	return nil
}
