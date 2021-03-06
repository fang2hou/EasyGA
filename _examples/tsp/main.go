package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/fang2hou/easyga"
	"github.com/wcharczuk/go-chart"
)

type travellingSalesmanProblem struct {
	ga           easyga.GeneticAlgorithm
	cityLocation [][]float64
	fitnessData  []float64
}

func main() {
	// Initialize a travelling salesman problem.
	var tsp travellingSalesmanProblem
	tsp.getCityLocation("tsp.cities.random.csv")
	tsp.DrawChart()
}

func (tsp *travellingSalesmanProblem) Init() {
	parameters := easyga.GeneticAlgorithmParameters{
		CrossoverProbability: .9,
		MutationProbability:  .35,
		PopulationSize:       2000,
		GenotypeNumber:       2,
		ChromosomeLength:     len(tsp.cityLocation),
		IterationsLimit:      30000,
		RandomSeed:           42,
		UseRoutine:           false,
	}

	custom := easyga.GeneticAlgorithmFunctions{
		ChromosomeInitFunction: func(c *easyga.Chromosome) {
			// Initialize
			c.Gene = make([]byte, 0)
			// Get an array contains the genes which tsp need.
			tspChromosome := easyga.Rand.Perm(parameters.ChromosomeLength)
			// Append each gene to the end of chromosome.
			for i := range tspChromosome {
				c.Gene = append(c.Gene, byte(tspChromosome[i]))
			}
		},
		MutateFunction: func(c *easyga.Chromosome) {
			// Get two different indexes of chromosome.
			index1 := c.GetRandomGeneIndex()
			index2 := c.GetRandomGeneIndex()
			for index1 == index2 {
				index2 = c.GetRandomGeneIndex()
			}
			// Switch value
			c.Gene[index1], c.Gene[index2] = c.Gene[index2], c.Gene[index1]
		},
		FitnessFunction: func(c *easyga.Chromosome) {
			// Initialize
			c.Fitness = 0

			// Be a travelling salesman :(
			for geneIndex := range c.Gene {
				// Get next city index from gene
				cityIndex := int(c.Gene[geneIndex])
				nextCityIndex := int(c.Gene[(geneIndex+1)%len(c.Gene)])
				// Calculate distance using pythagorean theorem
				distanceX := tsp.cityLocation[nextCityIndex][0] - tsp.cityLocation[cityIndex][0]
				distanceY := tsp.cityLocation[nextCityIndex][1] - tsp.cityLocation[cityIndex][1]
				distance := math.Sqrt(distanceX*distanceX + distanceY*distanceY)
				// Update fitness
				c.Fitness -= distance
			}
		},
		CrossOverFunction: func(parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
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
		},
		CheckStopFunction: func(ga *easyga.GeneticAlgorithm) bool {
			_, bestFitness := ga.Population.FindBest()
			maybeBest := float64(-5420)

			if bestFitness >= maybeBest || ga.Population.Iteration >= ga.Parameters.IterationsLimit {
				return true
			}
			return false
		},
		StatisticFunction: func(ga *easyga.GeneticAlgorithm) {
			_, bestFitness := ga.Population.FindBest()
			tsp.fitnessData = append(tsp.fitnessData, bestFitness)
		},
	}

	if err := tsp.ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}
}

func (tsp *travellingSalesmanProblem) Run() (easyga.Chromosome, float64, int) {
	return tsp.ga.Run()
}

func (tsp *travellingSalesmanProblem) getCityLocation(fileName string) {
	// Open file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create CSV Reader
	r := csv.NewReader(file)
	r.Comment = []rune("#")[0]

	// Parse data
	for {
		// Read a line before get EOF signal.
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// Add city location to array.
		tempCityX, _ := strconv.ParseFloat(record[0], 64)
		tempCityY, _ := strconv.ParseFloat(record[1], 64)
		tsp.cityLocation = append(tsp.cityLocation, []float64{tempCityX, tempCityY})
	}
}

func (tsp *travellingSalesmanProblem) DrawChart() {
	tsp.Init() // If you just want to run once, move tsp.init() to main().

	best, bestFit, iteration := tsp.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")

	xValue := make([]float64, 0)
	yValue := make([]float64, 0)

	for i := range best.Gene {
		xValue = append(xValue, float64(tsp.cityLocation[best.Gene[i]][0]))
		yValue = append(yValue, float64(tsp.cityLocation[best.Gene[i]][1]))
	}

	// Fix the line between the first city and last city.
	xValue = append(xValue, float64(tsp.cityLocation[best.Gene[0]][0]))
	yValue = append(yValue, float64(tsp.cityLocation[best.Gene[0]][1]))

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

	filePath := "stat.png"
	outFile, _ := os.Create(filePath)

	defer outFile.Close()

	graph.Render(chart.PNG, outFile)
}
