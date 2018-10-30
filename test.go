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

	//custom.ChromosomeInitFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.MutateFunction = func (parent1, parent2 *easyga.Chromosome) (child1, child2 *easyga.Chromosome) {
	//	You can customize your crossover function here
	//}

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
