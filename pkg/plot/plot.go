package plot

import (
	"io"

	"github.com/go-echarts/go-echarts/charts"
)

type Data struct {
	PlotName   string
	LabelXName string
	LabelYName string
	PlotLines  []Lines
}
type Lines struct {
	LineTitle string
	Points    []Point
}

type Point struct {
	X interface{}
	Y interface{}
}

// Draw plot data to out writer.
func Draw(data Data, out io.Writer) error {
	plot := charts.NewLine()
	plot.SetGlobalOptions(charts.TitleOpts{Title: data.PlotName})

	// p.X.Label.Text = data.LabelXName
	// p.Y.Label.Text = data.LabelXName

	for _, line := range data.PlotLines {
		xs, ys := line.toPlotterPoints()
		plot.AddXAxis(xs).
			AddYAxis(line.LineTitle, ys)
	}

	return plot.Render(out)
}

func (l Lines) toPlotterPoints() ([]interface{}, []interface{}) {
	count := len(l.Points)
	xs := make([]interface{}, count, count)
	ys := make([]interface{}, count, count)

	for index, point := range l.Points {
		xs[index] = point.X
		ys[index] = point.Y
	}

	return xs, ys
}
