package pipeviz

import (
	"fmt"
	svg "github.com/ajstarks/svgo/float"
	"io"
	"math"
)

func generateSvg(ir layouted, w io.Writer) {
	s := svg.New(w)
	var xSpacing float64 = 100
	var ySpacing float64 = 10

	totalWidth := float64(3000)
	totalHeight := float64(600)

	perColWidth := (totalWidth - float64(ir.Columns)*xSpacing) / float64(ir.Columns)
	perRowHeight := (totalHeight - float64(ir.Rows)*ySpacing) / float64(ir.Rows)

	s.Start(totalWidth, totalHeight)

	for _, b := range ir.Boxes {
		x := float64(b.StartCol) * (perColWidth + xSpacing)
		y := float64(b.StartRow) * (perRowHeight + ySpacing)
		width := perColWidth + float64(b.EndCol-b.StartCol)*(perColWidth+xSpacing)
		height := perRowHeight + float64(b.EndRow-b.StartRow)*(perRowHeight+ySpacing)

		stroke := `stroke="#ccc"`
		if b.Node.Class == "failed" {
			stroke = `stroke="#ff0000" stroke-width="2"`
		}
		s.Roundrect(x, y, width, height, 10, 10, `fill="#eee"`, stroke)

		rowPercentage := (float64(b.EndRow-b.StartRow) + 1) / float64(ir.Rows)

		if rowPercentage > 0.2 {
			fontSizeInt := int(math.Min(60, math.Floor(120*rowPercentage)))
			fontSize := fmt.Sprintf("font-size=\"%d\"", fontSizeInt)
			s.Text(x+(width/2)-150, y+(height/2)+8, b.Node.Label, fontSize, `dominant-baseline="middle"`)
		}

		lineStartX := x + width
		lineStartY := y + (height / 2)

		for _, lineRow := range b.Lines {
			endX := x + width + xSpacing
			endY := (lineRow * (perRowHeight + ySpacing)) + perRowHeight/2
			curviness := 0.6
			path := fmt.Sprintf("M %f %f C %f %f %f %f %f %f", lineStartX, lineStartY, lineStartX+(curviness*(endX-lineStartX)), lineStartY, endX-(curviness*(endX-lineStartX)), endY, endX, endY)
			s.Path(path, `fill="none" stroke="#aaa" stroke-width="4" stroke-dasharray="10"`)
		}
	}

	s.End()
}
