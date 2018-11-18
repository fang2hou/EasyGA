package easyga

// GeneticAlgorithmPopulation is a slice of Chromosome type array
type GeneticAlgorithmPopulation struct {
	Chromosomes []Chromosome
	Size        int
	Iteration   int
}

// GetRandomChromosomeIndex method will return a random chromosome index
func (gapop *GeneticAlgorithmPopulation) GetRandomChromosomeIndex() int {
	return Rand.Intn(gapop.Size)
}

// Init is a method to generate the first iteration
func (gapop *GeneticAlgorithmPopulation) Init(length int, size int, genotype int, initFunc func(c *Chromosome)) {
	// Save size
	gapop.Size = size

	// Create new individuals
	for i := 0; i < gapop.Size; i++ {
		// Initialize a new chromosome
		var tempChromosome Chromosome
		// Generate gene
		tempChromosome.GenotypeNumber = genotype
		if initFunc != nil {
			initFunc(&tempChromosome)
		} else {
			tempChromosome.Random(length, genotype)
		}
		// Add it into population
		gapop.Chromosomes = append(gapop.Chromosomes, tempChromosome)
	}
}

// FindBest is a method that return the index and fitness of the best one
func (gapop *GeneticAlgorithmPopulation) FindBest() (bestIndex int, bestFitness float64) {
	// Assume the first chromosome in population is the best
	bestIndex = 0
	bestFitness = gapop.Chromosomes[0].Fitness

	// If the chromosome better than the best one, set it as best
	for i := 1; i < gapop.Size; i++ {
		if gapop.Chromosomes[i].Fitness > bestFitness {
			bestIndex = i
			bestFitness = gapop.Chromosomes[i].Fitness
		}
	}

	return
}
