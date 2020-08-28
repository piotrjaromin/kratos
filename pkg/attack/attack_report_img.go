package attack

import (
	"fmt"
	"os"

	"github.com/piotrjaromin/kratos/pkg/plot"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type ImgReport struct {
	data plot.Data
	vegeta.Metrics
}

func NewImgReport() (*ImgReport, *vegeta.Metrics) {
	var m vegeta.Metrics
	plotLines := []plot.Lines{
		{
			LineTitle: "service",
			Points:    []plot.Point{},
		},
	}

	data := plot.Data{
		LabelXName: "time",
		LabelYName: "rps",
		PlotName:   "RPS",
		PlotLines:  plotLines,
	}

	return &ImgReport{
		data: data,
	}, &m
}

func (i *ImgReport) Add(res *vegeta.Result) {
	i.Metrics.Add(res)

	_, hour, min := i.End.Clock()

	point := plot.Point{
		X: fmt.Sprintf("%d:%d", hour, min),
		Y: i.Rate,
	}

	i.data.PlotLines[0].Points = append(i.data.PlotLines[0].Points, point)
}

func (i *ImgReport) Draw() error {
	file, err := os.Create("./plot.html")
	if err != nil {
		return err
	}

	return plot.Draw(i.data, file)
}
