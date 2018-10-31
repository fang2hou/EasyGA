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

	mapCityLocation := map[int][2]float64{
		0: {-1, 1}, 1: {0, 1}, 2: {1, 1},
		7: {-1, 0}, 3: {1, 0},
		6: {-1, -1}, 5: {0, -1}, 4: {1, -1},
	}

	parameters := easyga.Parameters{
		CrossoverProbability: .7,
		MutationProbability:  .1,
		PopulationSize:       10,
		Genotype:             2,
		ChromosomeLength:     8,
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
			nextCityIndex := int(c.Gene[(geneIndex + 1) % len(c.Gene)])
			// Calculate distance using pythagorean theorem
			distanceX := mapCityLocation[nextCityIndex][0] - mapCityLocation[cityIndex][0]
			distanceY := mapCityLocation[nextCityIndex][1] - mapCityLocation[cityIndex][1]
			distance := math.Sqrt(distanceX*distanceX + distanceY*distanceY)
			// Update fitness
			c.Fitness += distance
		}
	}

	custom.CrossOverFunction = func(parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
		length := len(parent1.Gene)
		// Init
		child1 = &easyga.Chromosome{Gene: make([]byte, 0)}
		child2 = &easyga.Chromosome{Gene: make([]byte, 0)}


		separatePoint1 := length/3
		if length % 3 != 0 {
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
		maybeBest := float64(8)

		if bestFitness <= maybeBest || ga.Iteration >= ga.Params.IterationsLimit {
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
