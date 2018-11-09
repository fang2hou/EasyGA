package EasyGA

import "math/rand"

type population struct {
	Chromosomes []Chromosome
}

func getRandomChromosomeIndex(p *population) int {
	return rand.Int() % len(p.Chromosomes)
}

func (p *population) Init(length int, size int, genotype uint8, initFunc func(c *Chromosome)) {
	p.Chromosomes = make([]Chromosome, 0)
	if initFunc != nil {
		for i := 0; i < size; i++ {
			var tempChromosome Chromosome
			initFunc(&tempChromosome)
			p.Chromosomes = append(p.Chromosomes, tempChromosome)
		}

		return
	}

	for i := 0; i < size; i++ {
		var tempChromosome Chromosome
		tempChromosome.Random(length, genotype)
		p.Chromosomes = append(p.Chromosomes, tempChromosome)
	}
}

func (p *population) FindBest() (bestIndex int, bestFitness float64) {
	bestIndex = 0
	bestFitness = p.Chromosomes[0].Fitness

	for i := range p.Chromosomes {
		if p.Chromosomes[i].Fitness > bestFitness {
			bestIndex = i
			bestFitness = p.Chromosomes[i].Fitness
		}
	}

	return
}
