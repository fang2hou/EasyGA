package main

import (
	"./easyga"
	"fmt"
	"math"
	"math/rand"
)

func main() {
	var ga easyga.GeneticAlgorithm

	rand.Seed(42)

	cityLocation := [8][2]float64{ {-1,1}, {0,1},{1,1},{-1,0},{1,0},{-1,-1},{0,-1},{1,-1}}[:]

	parameters := easyga.Parameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     10,
		IterationsLimit:      1000,
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

		// Cache the number of cities
		numberOfCities := len(cityLocation)

		// Be a travelling salesman :(
		for i, currentCityIndex := 0, 0; i < numberOfCities; i++{
			// Get next city index from gene
			nextCityIndex := int(c.Gene[currentCityIndex])

			// Calculate distance using pythagorean theorem
			distanceX := cityLocation[nextCityIndex][0] - cityLocation[currentCityIndex][0]
			distanceY := cityLocation[nextCityIndex][1] - cityLocation[currentCityIndex][1]
			distance := math.Sqrt(distanceX * distanceX + distanceY * distanceY)

			// Update fitness and currentCityIndex
			c.Fitness += distance
			currentCityIndex = nextCityIndex
		}
	}

	custom.CrossOverFunction = func (parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
		//Tsp
		if ga.Custom.CrossOverFunction != nil {
			return ga.Custom.CrossOverFunction(parent1, parent2)
		}

		// Default
		// - Single point crossover
		length := len(parent1.Gene)

		child1 = &easyga.Chromosome{Gene: make([]uint8, length)}
		child2 = &easyga.Chromosome{Gene: make([]uint8, length)}
		separatePoint1 := length / 3 + 1
		separatePoint2 := separatePoint1 * 2
		child2Center := parent1.Gene[separatePoint1:separatePoint2]
		child1Center := parent2.Gene[separatePoint1:separatePoint2]

		for i := 0;i < length; i++{
			isEqual := false
			for j := range child2Center{
				if parent1.Gene[i] == child2Center[j] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				child2.Gene[i] = parent1.Gene[i]
			}
		}

		tempchild2 := make([]byte, separatePoint2-separatePoint1+1)
		copy(tempchild2,child2.Gene[separatePoint1:separatePoint2])
		child2.Gene = append(child2.Gene[0:separatePoint1],child2Center...)
		child2.Gene = append(child2.Gene[:],tempchild2[:]...)

		for i := 0;i < length; i++{
			isEqual := false
			for j := range child2Center{
				if parent2.Gene[i] == child1Center[j] {
					isEqual = true
					break
				}
			}
			if !isEqual {
				child1.Gene[i] = parent2.Gene[i]
			}
		}

		tempchild1 := make([]byte, separatePoint2-separatePoint1+1)
		copy(tempchild1,child1.Gene[separatePoint1:separatePoint2])
		child1.Gene = append(child1.Gene[0:separatePoint1],child1Center...)
		child1.Gene = append(child1.Gene[:],tempchild1[:]...)


		return child1, child2
	}

	custom.CheckStopFunction = func (ga *easyga.GeneticAlgorithm) bool {
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
