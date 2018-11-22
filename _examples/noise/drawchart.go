package main

import (
	"os"

	"github.com/wcharczuk/go-chart"
)

func drawFitnessChart(fitnessArray []float64) {
	iterationLength := len(fitnessArray)
	xValues := make([]float64, 0)
	for i := 0; i < iterationLength; i++ {
		xValues = append(xValues, float64(i))

		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: xValues,
					YValues: fitnessArray,
				},
			},
		}
		filePath := "stat.png"
		outFile, _ := os.Create(filePath)

		defer outFile.Close()

		graph.Render(chart.PNG, outFile)

	}

}