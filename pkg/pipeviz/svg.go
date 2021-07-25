package pipeviz

import (
	"fmt"
	svg "github.com/ajstarks/svgo/float"
	"io"
)

func generateSvg(ir layouted, w io.Writer) {
	s := svg.New(w)
	var xSpacing float64 = 100
	var ySpacing float64 = 10

	totalWidth := float64(1800)
	totalHeight := float64(600)
	stepWidth := float64(300)
	stepHeight := float64(80)

	perColWidth := (totalWidth - float64(ir.Columns)*xSpacing) / float64(ir.Columns)
	perRowHeight := (totalHeight - float64(ir.Rows)*ySpacing) / float64(ir.Rows)

	s.Start(totalWidth, totalHeight)

	for _, b := range ir.Boxes {
		x := float64(b.StartCol) * (perColWidth + xSpacing)
		y := float64(b.StartRow) * (perRowHeight + ySpacing)
		width := perColWidth + float64(b.EndCol-b.StartCol)*(perColWidth+xSpacing)
		height := perRowHeight + float64(b.EndRow-b.StartRow)*(perRowHeight+ySpacing)

		halfH := height/2 - stepHeight/2

		stroke := `stroke="#ccc"`
		if b.Node.Class == "failed" {
			stroke = `stroke="#ff0000" stroke-width="2"`
		}
		//s.Roundrect(x, y, width, height, 10, 10, `fill="pink"`)
		s.Roundrect(x, y+halfH, stepWidth, stepHeight, 40, 40, `fill="#eee"`,stroke)
		s.Circle(x+40, y+ (height/2), 32, `fill="#fff"`,`stroke="#ccc"`)

		if b.Node.Class == "succeeded" {
			translate := fmt.Sprintf(`transform="translate(%f, %f)"`, x+15, y+halfH+15)
			s.Polygon([]float64{40.6, 17, 7.4, 4.6, 17, 43.4}, []float64{12.1, 35.7,26.1,29,41.3,14.9}, `fill="#43A047"`, translate)
		} else if b.Node.Class == "failed" {
			translate := fmt.Sprintf(`transform="translate(%f, %f),scale(0.3)"`, x+25, y+halfH+25)
			s.Group(translate)
			s.Path("M 6.3895625,6.4195626 C 93.580437,93.610437 93.580437,93.610437 93.580437,93.610437", `style="fill:none;fill-rule:evenodd;stroke:#ff0000;stroke-width:18.05195999;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1"`)
			s.Path("M 6.3894001,93.6106 C 93.830213,6.4194003 93.830213,6.4194003 93.830213,6.4194003", `fill:none;fill-rule:evenodd;stroke:#ff0000;stroke-width:17.80202103;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1`)
			s.Gend()
		} else if b.Node.Class == "running" {
			transform := fmt.Sprintf(`transform="translate(%f, %f),scale(0.025)"`, x+21, y+halfH+24)
			s.Path("m 1024,640 q 0,106 -75,181 -75,75 -181,75 -106,0 -181,-75 -75,-75 -75,-181 0,-106 75,-181 75,-75 181,-75 106,0 181,75 75,75 75,181 z m 512,109 V 527 q 0,-12 -8,-23 -8,-11 -20,-13 l -185,-28 q -19,-54 -39,-91 35,-50 107,-138 10,-12 10,-25 0,-13 -9,-23 -27,-37 -99,-108 -72,-71 -94,-71 -12,0 -26,9 l -138,108 q -44,-23 -91,-38 -16,-136 -29,-186 -7,-28 -36,-28 H 657 q -14,0 -24.5,8.5 Q 622,-111 621,-98 L 593,86 q -49,16 -90,37 L 362,16 Q 352,7 337,7 323,7 312,18 186,132 147,186 q -7,10 -7,23 0,12 8,23 15,21 51,66.5 36,45.5 54,70.5 -27,50 -41,99 L 29,495 Q 16,497 8,507.5 0,518 0,531 v 222 q 0,12 8,23 8,11 19,13 l 186,28 q 14,46 39,92 -40,57 -107,138 -10,12 -10,24 0,10 9,23 26,36 98.5,107.5 72.5,71.5 94.5,71.5 13,0 26,-10 l 138,-107 q 44,23 91,38 16,136 29,186 7,28 36,28 h 222 q 14,0 24.5,-8.5 Q 914,1391 915,1378 l 28,-184 q 49,-16 90,-37 l 142,107 q 9,9 24,9 13,0 25,-10 129,-119 165,-170 7,-8 7,-22 0,-12 -8,-23 -15,-21 -51,-66.5 -36,-45.5 -54,-70.5 26,-50 41,-98 l 183,-28 q 13,-2 21,-12.5 8,-10.5 8,-23.5 z", `fill="#89CFF0"`, transform)
		}

		s.Text(x+100, y+(height/2)+8, b.Node.Label, `font-size="24px"`, `dominant-baseline="middle"`)

		lineStartX := x+width - 20
		lineStartY := y+(height/2)

		if b.Lines != nil {
			s.Line(x+stepWidth, y+(height/2), lineStartX, lineStartY, `stroke="#aaa" stroke-width="4" stroke-dasharray="10"`)
		}
		for _, lineRow := range b.Lines {
			endX := x+width+ xSpacing
			endY := (lineRow * (perRowHeight+ ySpacing)) + perRowHeight/2
			curviness := 0.6
			path := fmt.Sprintf("M %f %f C %f %f %f %f %f %f", lineStartX, lineStartY, lineStartX + (curviness*(endX-lineStartX)), lineStartY, endX-(curviness*(endX-lineStartX)), endY, endX, endY)
			s.Path(path, `fill="none" stroke="#aaa" stroke-width="4" stroke-dasharray="10"`)
		}
	}

	s.End()
}
