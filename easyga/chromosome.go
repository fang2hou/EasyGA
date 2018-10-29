package easyga

import (
	"math/rand"
)

type Chromosome struct {
	Gene    []byte
	Fitness float64
}

func getRandomGenotype(genotype uint8) byte {
	return byte(rand.Intn(int(genotype)))
}

func (c *Chromosome) Length() int {
	return len(c.Gene)
}

func (c *Chromosome) GetRandomGeneIndex() int {
	return rand.Int() % len(c.Gene)
}

func (c *Chromosome) Random(length int, genotype uint8) {
	c.Gene = make([]byte, 0)
	for i := 0; i < length; i++ {
		c.Gene = append(c.Gene, getRandomGenotype(genotype))
	}
}

func (c *Chromosome) UpdateFitness() {
	// TODO: Make Fitness function be customizable
	c.Fitness = 0
	for _, genotype := range c.Gene {
		c.Fitness += float64(genotype)
	}
}

func (c *Chromosome) Mutate(genotype uint8) {
	c.Gene[c.GetRandomGeneIndex()] = getRandomGenotype(genotype)
	c.UpdateFitness()
}
