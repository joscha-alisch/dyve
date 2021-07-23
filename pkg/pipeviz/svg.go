package pipeviz

import (
	svg "github.com/ajstarks/svgo/float"
	"io"
)

func generateSvg(ir layouted, w io.Writer) {
	s := svg.New(w)
	var spacing float64 = 20
	totalWidth := float64(1800)
	totalHeight := float64(600)

	perColWidth := (totalWidth - float64(ir.Columns)*spacing) / float64(ir.Columns)
	perRowHeight := (totalHeight - float64(ir.Rows)*spacing) / float64(ir.Rows)

	s.Start(totalWidth, totalHeight)

	for _, b := range ir.Boxes {
		x := float64(b.StartCol) * (perColWidth + spacing)
		y := float64(b.StartRow) * (perRowHeight + spacing)

		width := perColWidth + float64(b.EndCol-b.StartCol)*(perColWidth+spacing)
		height := perRowHeight + float64(b.EndRow-b.StartRow)*(perRowHeight+spacing)

		s.Roundrect(x, y, width, height, 10, 10, `fill="white"`)
		s.Text(x+(width/2), y+(height/2), b.Node.Label, `font-size="24px"`, `dominant-baseline="middle"`, `text-anchor="middle"`)
	}

	s.End()
}
