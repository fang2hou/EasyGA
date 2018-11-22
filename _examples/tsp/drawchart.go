package main

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func drawFitnessChart(fitnessArray []float64) {
	iterationLength := len(fitnessArray)
	xValues := make([]float64, 0)
	for i := 0; i < iterationLength; i++ {
		xValues = append(xValues, float64(i))
	}
	fitnessSeries := chart.ContinuousSeries{
		XValues: xValues,
		YValues: fitnessArray,
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorBlack,
		},
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.0f", vf)
				}
				return ""
			},
		},

		YAxis: chart.YAxis{
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprintf("%0.0f", vf)
				}
				return ""
			},
			AxisType: chart.YAxisSecondary,
		},

		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},

		Series: []chart.Series{
			fitnessSeries,
		},
	}
	filePath := "stat.png"
	outFile, _ := os.Create(filePath)

	defer outFile.Close()

	graph.Render(chart.PNG, outFile)

}
