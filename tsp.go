package main

import (
	"./easyga"
	"fmt"
	"math/rand"
)

func main() {
	var ga easyga.GeneticAlgorithm

	rand.Seed(42)

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
		c.Gene = make([]byte, 0)

		tspChromosome := rand.Perm(parameters.ChromosomeLength)

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

	//custom.FitnessFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.CrossOverFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.CheckStopFunction = func (ga *easyga.GeneticAlgorithm) bool {
	//	You can customize your check stop function here
	//}

	if err := ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}

	best, bestFit, iteration := ga.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")
}
