package easyga

import "math/rand"

type population struct {
	chromosomes []Chromosome
}

func getRandomChromosomeIndex(p *population) int {
	return rand.Int() % len(p.chromosomes)
}

func (p *population) Init(length int, size int, genotype uint8, initFunc func(c *Chromosome)) {
	p.chromosomes = make([]Chromosome, 0)
	if initFunc != nil {
		for i := 0; i < size; i++ {
			var tempChromosome Chromosome
			initFunc(&tempChromosome)
			p.chromosomes = append(p.chromosomes, tempChromosome)
		}

		return
	}

	for i := 0; i < size; i++ {
		var tempChromosome Chromosome
		tempChromosome.Random(length, genotype)
		p.chromosomes = append(p.chromosomes, tempChromosome)
	}
}

func (p *population) FindBest() (bestIndex int, bestFitness float64) {
	bestIndex = 0
	bestFitness = p.chromosomes[0].Fitness

	for i := range p.chromosomes {
		if p.chromosomes[i].Fitness > bestFitness {
			bestIndex = i
			bestFitness = p.chromosomes[i].Fitness
		}
	}

	return
}
