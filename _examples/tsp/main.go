package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/fang2hou/easyga"
	"github.com/wcharczuk/go-chart"
)

func main() {
	easyga.
	fmt.Println("Start service at http://localhost:8182/")
	http.HandleFunc("/", drawChart)
	http.ListenAndServe(":8182", nil)
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	var ga easyga.GeneticAlgorithm

	cityLocation := readCSVFile()

	parameters := easyga.Parameters{
		CrossoverProbability: .8,
		MutationProbability:  .2,
		PopulationSize:       20,
		Genotype:             2,
		ChromosomeLength:     len(cityLocation),
		IterationsLimit:      100000,
	}

	custom := easyga.CustomFunctions{}

	custom.ChromosomeInitFunction = func(c *easyga.Chromosome) {
		// Initialize
		c.Gene = make([]byte, 0)
		// Get a array contains the genes which tsp need
		tspChromosome := rand.Perm(parameters.ChromosomeLength)
		// Append each gene to chromosome
		for i := range tspChromosome {
			c.Gene = append(c.Gene, byte(tspChromosome[i]))
		}
	}

	custom.MutateFunction = func(c *easyga.Chromosome) {
		// Get two different index of chromosome
		index1 := c.GetRandomGeneIndex()
		index2 := c.GetRandomGeneIndex()
		for index1 == index2 {
			index2 = c.GetRandomGeneIndex()
		}
		// Switch value
		c.Gene[index1], c.Gene[index2] = c.Gene[index2], c.Gene[index1]
	}

	custom.FitnessFunction = func(c *easyga.Chromosome) {
		// Initialize
		c.Fitness = 0

		// Be a travelling salesman :(
		for geneIndex := range c.Gene {
			// Get next city index from gene
			cityIndex := int(c.Gene[geneIndex])
			nextCityIndex := int(c.Gene[(geneIndex+1)%len(c.Gene)])
			// Calculate distance using pythagorean theorem
			distanceX := cityLocation[nextCityIndex][0] - cityLocation[cityIndex][0]
			distanceY := cityLocation[nextCityIndex][1] - cityLocation[cityIndex][1]
			distance := math.Sqrt(distanceX*distanceX + distanceY*distanceY)
			// Update fitness
			c.Fitness -= distance
		}
	}

	custom.CrossOverFunction = func(parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
		// Find separate part
		crossoverStart := parameters.ChromosomeLength / 3
		if parameters.ChromosomeLength%3 != 0 {
			crossoverStart++
		}
		crossoverEnd := crossoverStart * 2

		// child 1
		child1 = &easyga.Chromosome{Gene: make([]byte, 0)}
		crossoverPart := parent2.Gene[crossoverStart:crossoverEnd]
		copyPart := make([]byte, 0)
		for parentIndex := range parent1.Gene {
			isEqual := false
			for skipCopyIndex := range crossoverPart {
				if parent1.Gene[parentIndex] == crossoverPart[skipCopyIndex] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				copyPart = append(copyPart, parent1.Gene[parentIndex])
			}
		}

		child1.Gene = append(child1.Gene, copyPart[0:crossoverStart]...)
		child1.Gene = append(child1.Gene, crossoverPart...)
		child1.Gene = append(child1.Gene, copyPart[crossoverStart:]...)

		// child 2
		child2 = &easyga.Chromosome{Gene: make([]byte, 0)}
		crossoverPart = parent1.Gene[crossoverStart:crossoverEnd]
		copyPart = make([]byte, 0)
		for parentIndex := range parent2.Gene {
			isEqual := false
			for skipCopyIndex := range crossoverPart {
				if parent2.Gene[parentIndex] == crossoverPart[skipCopyIndex] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				copyPart = append(copyPart, parent2.Gene[parentIndex])
			}
		}

		child2.Gene = append(child2.Gene, copyPart[0:crossoverStart]...)
		child2.Gene = append(child2.Gene, crossoverPart...)
		child2.Gene = append(child2.Gene, copyPart[crossoverStart:]...)

		return
	}

	custom.CheckStopFunction = func(ga *easyga.GeneticAlgorithm) bool {
		_, bestFitness := ga.Population.FindBest()
		maybeBest := float64(-1877.214)

		if bestFitness >= maybeBest || ga.Iteration >= ga.Params.IterationsLimit {
			return true
		}

		return false
	}

	if err := ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}

	best, bestFit, iteration := ga.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")

	xValue := make([]float64, 0)
	yValue := make([]float64, 0)

	for i := range best.Gene {
		xValue = append(xValue, float64(cityLocation[best.Gene[i]][0]))
		yValue = append(yValue, float64(cityLocation[best.Gene[i]][1]))
	}

	// Fix the line between the first city and last city
	xValue = append(xValue, float64(cityLocation[best.Gene[0]][0]))
	yValue = append(yValue, float64(cityLocation[best.Gene[0]][1]))

	tspSeries := chart.ContinuousSeries{
		XValues: xValue,
		YValues: yValue,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlack,
			StrokeWidth: 1.0,
		},
	}

	graph := chart.Chart{
		Title:  "TSP Final result",
		Width:  500,
		Height: 500,
		DPI:    100.0,
		Series: []chart.Series{tspSeries},
	}

	res.Header().Set("Content-Type", chart.ContentTypePNG)
	err := graph.Render(chart.PNG, res)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func readCSVFile() (cityLocation [][]float64) {
	fileName := "tsp.cities.cycle.csv"
	ioReader, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(strings.NewReader(string(ioReader)))

	// Skip the line start with #
	r.Comment = []rune("#")[0]

	// Parse data
	for i := 0; ; i++ {
		record, err := r.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		tempCityX, _ := strconv.ParseFloat(record[0], 64)
		tempCityY, _ := strconv.ParseFloat(record[1], 64)
		cityLocation = append(cityLocation, []float64{tempCityX, tempCityY})
	}

	return
}
