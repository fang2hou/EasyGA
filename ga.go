package easyga

import (
	"math/rand"
	"sync"
)

// Rand is a pointer to the rand which can be changed outside
var Rand *rand.Rand

// GeneticAlgorithm is a struct that contains everything of genetic algorithm.
type GeneticAlgorithm struct {
	Parameters GeneticAlgorithmParameters
	Functions  GeneticAlgorithmFunctions
	Population GeneticAlgorithmPopulation
}

// Init method will initialize the original population.
func (ga *GeneticAlgorithm) Init(customParameters GeneticAlgorithmParameters, customFunctions GeneticAlgorithmFunctions) (err error) {
	// Check parameters before initialization
	if err = customParameters.check(); err == nil {
		// Initialize parameters
		ga.Parameters = customParameters

		// Initialize Seed of rand
		Rand = rand.New(rand.NewSource(ga.Parameters.RandomSeed))

		// Initialize functions
		ga.Functions = customFunctions
		ga.Functions.Init()

		// Initialize population
		ga.Population.Iteration = 0
		ga.Population.Init(ga.Parameters.ChromosomeLength, ga.Parameters.PopulationSize, ga.Parameters.GenotypeNumber, ga.Functions.ChromosomeInitFunction)

		// Update fitness of first generation
		ga.updateFitness()
	}

	return
}

// Run method will create a loop for find best result.
func (ga *GeneticAlgorithm) Run() (best Chromosome, fitness float64, iteration int) {
	for !ga.checkStop() {
		// Selection - Select parents from population
		parentsPair := ga.selection()

		// Crossover - perform crossover on parents creating population
		ga.crossover(parentsPair)

		// Mutation - perform mutation of population
		ga.mutation()

		// Update fitness
		ga.updateFitness()

		// Update iteration
		ga.Population.Iteration++
	}

	bestIndex, bestFitness := ga.Population.FindBest()
	fitness = bestFitness
	iteration = ga.Population.Iteration

	best = ga.Population.Chromosomes[bestIndex]

	return
}

func (ga *GeneticAlgorithm) selection() (parentsPair [][2]int) {
	selectedPopulation := ga.Functions.SelectFunction(ga)
	selector := 0

	for selector < ga.Population.Size {
		parentIndex1, parentIndex2 := selectedPopulation[selector], selectedPopulation[selector+1]
		selector += 2

		parentsPair = append(parentsPair, [2]int{parentIndex1, parentIndex2})
	}

	return parentsPair
}

func (ga *GeneticAlgorithm) crossover(parentsPair [][2]int) {
	var nextPopulation GeneticAlgorithmPopulation

	// Copy information to next population
	nextPopulation.Size = ga.Population.Size
	nextPopulation.Iteration = ga.Population.Iteration

	for i := 0; i < len(parentsPair); i++ {
		// Get the indexes of parents
		parent1 := ga.Population.Chromosomes[parentsPair[i][0]]
		parent2 := ga.Population.Chromosomes[parentsPair[i][1]]

		// Initialize children chromosome
		var child1, child2 *Chromosome

		// Crossover with probability
		if Rand.Float64() < ga.Parameters.CrossoverProbability {
			child1, child2 = ga.Functions.CrossOverFunction(&parent1, &parent2)
			child1.GenotypeNumber = ga.Parameters.GenotypeNumber
			child2.GenotypeNumber = ga.Parameters.GenotypeNumber
		} else {
			child1, child2 = &parent1, &parent2
		}

		// Add child to the next generation population
		nextPopulation.Chromosomes = append(nextPopulation.Chromosomes, *child1, *child2)
	}

	// Update population
	ga.Population = nextPopulation
}

func (ga *GeneticAlgorithm) mutation() {
	var routineWait sync.WaitGroup

	for i := 0; i < ga.Population.Size; i++ {
		// Mutate with probability
		if Rand.Float64() < ga.Parameters.MutationProbability {
			routineWait.Add(1)

			go func(i int, counter *sync.WaitGroup) {
				ga.Functions.MutateFunction(&ga.Population.Chromosomes[i])
				counter.Done()
			}(i, &routineWait)
		}
	}

	routineWait.Wait()
}

func (ga *GeneticAlgorithm) updateFitness() {
	var routineWait sync.WaitGroup

	for i := 0; i < ga.Population.Size; i++ {
		routineWait.Add(1)

		go func(index int, counter *sync.WaitGroup) {
			ga.Functions.FitnessFunction(&ga.Population.Chromosomes[index])
			counter.Done()
		}(i, &routineWait)
	}

	routineWait.Wait()
}

func (ga *GeneticAlgorithm) checkStop() bool {
	return ga.Functions.CheckStopFunction(ga)
}
