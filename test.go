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

	//ga.CheckStopFunction = func (ga *easyga.GeneticAlgorithm) bool {
	//	_, bestFitness := ga.Population.FindBest()
	//	maybeBest := int(ga.Params.Genotype-1) * ga.Params.ChromosomeLength - 1
	//
	//	if int(bestFitness) >= maybeBest || ga.Iteration >= ga.Params.IterationsLimit {
	//		return true
	//	}
	//
	//	return false
	//}

	if err := ga.Init(parameters); err != nil {
		fmt.Println(err)
		return
	}

	best, bestFit, iteration := ga.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")
}
