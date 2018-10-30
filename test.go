package main

import (
	"./easyga"
	"fmt"
	"math"
	"math/rand"
)

func main() {
	var ga easyga.GeneticAlgorithm

	parameters := easyga.Parameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     10,
		IterationsLimit:      1000,
	}

	custom := easyga.CustomFunctions{}

	//custom.CheckStopFunction = func (ga *easyga.GeneticAlgorithm) bool {
	//	You can customize your check stop function here
	//}

	custom.CrossOverFunction = func (parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
		//Tsp
		if ga.Custom.CrossOverFunction != nil {
			return ga.Custom.CrossOverFunction(parent1, parent2)
		}

		// Default
		// - Single point crossover
		length := len(parent1.Gene)
		position := parent1.GetRandomGeneIndex()

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
		child2.Gene = append(child2.Gene[0:separatePoint1],child2Center...,tempchild2[:]...)

		return child1, child2
	}
	//
	custom.FitnessFunction = func(c *easyga.Chromosome) {
		//Tsp
		c.Fitness = 0
		for _, genotype := range c.Gene {
			c.Fitness += float64(genotype)
			const dimension int = 2
			const placeNumber int = 2
			location := [placeNumber][dimension]float64{}

			for i := 0;i < placeNumber; i++{
				for j := 0 ; j < dimension ; j++{
					location[i][j] = float64(rand.Int()% 10)
				}
			}


			for i := 0; i < placeNumber- 1; i++ {
				genotype := c.Gene[i]
				xDistance := location[genotype][0] - location[genotype+1][0]
				yDistance := location[genotype][1] - location[genotype+1][1]
				distance := math.Sqrt(xDistance * xDistance + yDistance * yDistance)
				c.Fitness += distance
			}
		}
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
