package easyga

import "math/rand"

type population struct {
	chromosomes []Chromosome
}

func getRandomChromosomeIndex(p *population) int {
	return rand.Int() % len(p.chromosomes)
}

func (p *population) Init(length int, size int, genotype uint8) {
	p.chromosomes = make([]Chromosome, 0)
	for i := 0; i < size; i++ {
		var tempChromosome Chromosome
		tempChromosome.Random(length, genotype)
		tempChromosome.UpdateFitness()
		p.chromosomes = append(p.chromosomes, tempChromosome)
	}
}

func (p *population) FindBest() (bestIndex int, bestFitness float64) {
	bestIndex = 0
	bestFitness = p.chromosomes[0].fitness

	for i := range p.chromosomes {
		if p.chromosomes[i].fitness > bestFitness {
			bestIndex = i
			bestFitness = p.chromosomes[i].fitness
		}
	}

	return
}
