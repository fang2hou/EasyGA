package easyga

import "math/rand"

type population struct {
	Chromosomes []Chromosome
}

// getRandomChromosomeIndex is a private method that return a random chromosome index
func getRandomChromosomeIndex(p *population) int {
	return rand.Int() % len(p.Chromosomes)
}

// Init is a method to generate the first iteration
func (p *population) Init(length int, size int, genotype uint8, initFunc func(c *Chromosome)) {
	// Initialize population
	p.Chromosomes = make([]Chromosome, 0)

	// Create new individuals
	for i := 0; i < size; i++ {
		// Initialize a new chromosome
		var tempChromosome Chromosome
		// Generate gene
		if initFunc != nil {
			initFunc(&tempChromosome)
		} else {
			tempChromosome.Random(length, genotype)
		}
		// Add it into population
		p.Chromosomes = append(p.Chromosomes, tempChromosome)
	}
}

// FindBest is a method that return the index and fitness of the best one
func (p *population) FindBest() (bestIndex int, bestFitness float64) {
	// Assume the first chromosome in population is the best
	bestIndex = 0
	bestFitness = p.Chromosomes[0].Fitness

	// If the chromosome better than the best one, set it as best
	for i := range p.Chromosomes {
		if p.Chromosomes[i].Fitness > bestFitness {
			bestIndex = i
			bestFitness = p.Chromosomes[i].Fitness
		}
	}

	return
}
