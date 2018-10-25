package main

import (
	"fmt"

	"./easyga"
)

func main() {
	var ga easyga.GeneticAlgorithm

	parameters := easyga.Parameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     10,
		IterationsLimit:      100000,
	}

	if err := ga.Init(parameters); err != nil {
		fmt.Println(err)
		return
	}

	best, bestFit, iteration := ga.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")
}
