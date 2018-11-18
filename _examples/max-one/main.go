package main

import (
	"fmt"

	"github.com/fang2hou/easyga"
)

func main() {
	var ga easyga.GeneticAlgorithm

	parameters := easyga.GeneticAlgorithmParameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     10,
		IterationsLimit:      1000,
		RandomSeed:           43,
	}

	custom := easyga.GeneticAlgorithmFunctions{}

	//custom.ChromosomeInitFunction = func(c *easyga.Chromosome) {
	//	You can customize your chromosome initialization function here
	//}

	//custom.SelectFunction =  func(ga *GeneticAlgorithm) []int {
	//	You can customize your selection function here
	//}

	//custom.CrossOverFunction = func(c *easyga.Chromosome) {
	//	You can customize your crossover function here
	//}

	//custom.MutateFunction = func(c *easyga.Chromosome) {
	//	You can customize your mutation function here
	//}

	//custom.FitnessFunction = func(c *easyga.Chromosome) {
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
