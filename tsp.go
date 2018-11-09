package main

import (
	"./easyga"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func main() {
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
		length := len(parent1.Gene)
		// Init
		child1 = &easyga.Chromosome{Gene: make([]byte, 0)}
		child2 = &easyga.Chromosome{Gene: make([]byte, 0)}

		separatePoint1 := length / 3
		if length%3 != 0 {
			separatePoint1 += 1
		}
		separatePoint2 := separatePoint1 * 2
		// Child1
		child1Center := parent2.Gene[separatePoint1:separatePoint2]
		tempChild1Gene := make([]byte, 0)
		for i := range parent1.Gene {
			isEqual := false
			for j := range child1Center {
				if parent1.Gene[i] == child1Center[j] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				tempChild1Gene = append(tempChild1Gene, parent1.Gene[i])
			}
		}

		child1.Gene = append(child1.Gene, tempChild1Gene[0:separatePoint1]...)
		child1.Gene = append(child1.Gene, child1Center...)
		child1.Gene = append(child1.Gene, tempChild1Gene[separatePoint1:]...)
		// Child2
		child2Center := parent1.Gene[separatePoint1:separatePoint2]
		tempChild2Gene := make([]byte, 0)
		for i := range parent2.Gene {
			isEqual := false
			for j := range child2Center {
				if parent2.Gene[i] == child2Center[j] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				tempChild2Gene = append(tempChild2Gene, parent2.Gene[i])
			}
		}

		child2.Gene = append(child2.Gene, tempChild2Gene[0:separatePoint1]...)
		child2.Gene = append(child2.Gene, child2Center...)
		child2.Gene = append(child2.Gene, tempChild2Gene[separatePoint1:]...)

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
}

func readCSVFile() [][]float64 {
	cityLocation := make([][]float64, 0)

	fileName := "./tsp.cities.csv"
	ioReader, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	r := csv.NewReader(strings.NewReader(string(ioReader)))

	// Disable the sentence start with #
	r.Comment = []rune("#")[0]

	for i := 0; ; i++ {

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		tempCityX, _ := strconv.ParseFloat(record[0], 64)
		tempCityY, _ := strconv.ParseFloat(record[1], 64)

		cityLocation = append(cityLocation, []float64{tempCityX, tempCityY})

	}

	return cityLocation
}
