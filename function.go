package easyga

// GeneticAlgorithmFunctions is a struct that contains every functions customized.
type GeneticAlgorithmFunctions struct {
	ChromosomeInitFunction func(c *Chromosome)

	SelectFunction    func(ga *GeneticAlgorithm) []int
	CrossOverFunction func(parent1, parent2 *Chromosome) (child1, child2 *Chromosome)
	MutateFunction    func(c *Chromosome)
	FitnessFunction   func(c *Chromosome)

	CheckStopFunction func(ga *GeneticAlgorithm) bool
}

// Init method will generate functions when not be initialized.
func (gafuncs *GeneticAlgorithmFunctions) Init() {
	// Mutation
	if gafuncs.MutateFunction == nil {
		gafuncs.MutateFunction = func(c *Chromosome) {
			// Replace a genotype with a new one generated randomly.
			c.Gene[c.GetRandomGeneIndex()] = c.GetRandomGenotype()
		}
	}

	// Fitness
	if gafuncs.FitnessFunction == nil {
		gafuncs.FitnessFunction = func(c *Chromosome) {
			// The sum of the number of genotype
			c.Fitness = 0
			for i := range c.Gene {
				c.Fitness += float64(c.Gene[i])
			}
		}
	}

	// Select
	if gafuncs.SelectFunction == nil {
		gafuncs.SelectFunction = TournamentFunction
	}

	// Crossover
	if gafuncs.CrossOverFunction == nil {
		gafuncs.CrossOverFunction = func(parent1, parent2 *Chromosome) (child1, child2 *Chromosome) {
			// Single-point crossover
			length := len(parent1.Gene)
			position := parent1.GetRandomGeneIndex()

			child1 = &Chromosome{Gene: make([]byte, length)}
			child2 = &Chromosome{Gene: make([]byte, length)}

			copy(child1.Gene, parent1.Gene[:position])
			copy(child2.Gene, parent2.Gene[:position])
			child1.Gene = append(child1.Gene[:position], parent2.Gene[position:]...)
			child2.Gene = append(child2.Gene[:position], parent1.Gene[position:]...)

			return child1, child2
		}
	}

	// Check stop
	if gafuncs.CheckStopFunction == nil {
		gafuncs.CheckStopFunction = func(ga *GeneticAlgorithm) bool {
			// The maximum sum of the number of genotype
			_, bestFitness := ga.Population.FindBest()
			maybeBest := int(ga.Parameters.Genotype-1) * int(ga.Parameters.ChromosomeLength)

			if int(bestFitness) >= maybeBest || ga.Population.Iteration >= ga.Parameters.IterationsLimit {
				return true
			}

			return false
		}
	}
}

// TournamentFunction is the default select method in EasyGA
func TournamentFunction(ga *GeneticAlgorithm) (newPopulationIndex []int) {
	// Competition function
	findWinner := func(index1 int, index2 int) (winner int) {
		if winner = index1; ga.Population.Chromosomes[index2].Fitness > ga.Population.Chromosomes[index1].Fitness {
			winner = index2
		}

		return
	}

	// Tournament!
	for i := 0; i < ga.Parameters.PopulationSize; i++ {
		chromosomeIndex1 := ga.Population.GetRandomChromosomeIndex()
		chromosomeIndex2 := ga.Population.GetRandomChromosomeIndex()

		newPopulationIndex = append(newPopulationIndex, findWinner(chromosomeIndex1, chromosomeIndex2))
	}

	return
}
