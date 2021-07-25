package pipeviz

type grid struct {
	constraints      []constraint
	constraintSource map[int][]constraint
	constraintTarget map[int][]constraint

	width  int
	height int
	rows   map[int]map[int]bool
	cols   map[int]map[int]bool
	boxes  map[int]*position
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
		boxes:            make(map[int]*position),
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
	if g.width == 0 {
		g.width = 1
	}
	g.height++

	p := &position{
		row:    g.height-1,
		col:    0,
		width:  1,
		height: 1,
	}

	g.boxes[id] = p

	for i := p.row; i < p.row+p.height; i++ {
		if g.rows[i] == nil {
			g.rows[i] = make(map[int]bool)
		}
		g.rows[i][id] = true
	}
	for i := p.col; i < p.col+p.width; i++ {
		if g.cols[i] == nil {
			g.cols[i] = make(map[int]bool)
		}
		g.cols[i][id] = true
	}
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

	targetCol := g.boxes[target].col
	srcCol := g.boxes[src].col

	if srcCol > targetCol {
		return
	}

	g.boxes[src].col = targetCol+1
	if g.cols[targetCol+1] == nil {
		g.cols[targetCol+1] = make(map[int]bool)
	}

	g.cols[targetCol+1][src] = true

	if g.width <= targetCol+1 {
		g.width++
	}

	for i := srcCol; i <= targetCol; i++ {
		delete(g.cols[i], src)
	}
}

func (g *grid) layout() {
	for id := range g.boxes {
		maxY := g.width
		for _, constraint := range g.constraintTarget[id] {
			if g.boxes[constraint.src].col < maxY {
				maxY = g.boxes[constraint.src].col
			}
		}
		g.boxes[id].width = maxY - g.boxes[id].col

		for col := g.boxes[id].col; col < g.boxes[id].col + g.boxes[id].width ; col++ {
			g.cols[col][id] = true
		}
	}

	maxHeight := 0
	for i := 0; i < g.width; i++ {
		g.reorderColumn(i)
		if maxHeight < len(g.cols[i]) {
			maxHeight = len(g.cols[i])
		}
	}

	g.height = maxHeight

	for col := range g.cols {
		resizable := make([]*int, g.height)
		for id := range g.cols[col] {
			id := id
			if g.boxes[id].width == 1 {
				resizable[g.boxes[id].row] = &id
			}
		}

		space := g.height - len(g.cols[col])
		i := 0
		for i < space {
			added := 0
			for _, id := range resizable {
				if id == nil {
					continue
				}

				if i >= space {
					g.boxes[*id].row += added
					continue
				}
				g.boxes[*id].row += added
				g.boxes[*id].height += 1
				i++
				added++
			}
		}
	}
}

func (g *grid) reorderColumn(col int) {
	single := make(map[int]bool)
	multiStarting := make(map[int]bool)
	multi := make(map[int]bool)

	for id, b := range g.cols[col] {
		if !b {
			continue
		}

		if g.boxes[id].col == col {
			if g.boxes[id].width == 1 {
				single[id] = true
			} else {
				multiStarting[id] = true
			}
		} else {
			multi[id] = true
		}
	}

	row := 0
	for id := range single {
		g.moveToRow(id, row)
		row += g.boxes[id].height
	}

	for id := range multiStarting {
		g.moveToRow(id, row)
		row += g.boxes[id].height
	}

	for id := range multi {
		g.moveToRow(id, row)
		row += g.boxes[id].height
	}
}

func (g *grid) moveToRow(id int, row int) {
	before := g.boxes[id].row
	g.boxes[id].row = row
	delete(g.rows[before], id)
	g.rows[row][id] = true
}

func (g *grid) getBoxes() map[int]box {
	b := make(map[int]box)

	for id, p := range g.boxes {
		b[id] = box{
			StartCol: p.col,
			EndCol:   p.col + (p.width - 1),
			StartRow: p.row,
			EndRow:   p.row + (p.height - 1),
			Node: Node{
				Id: id,
			},
		}
	}
	return b
}
