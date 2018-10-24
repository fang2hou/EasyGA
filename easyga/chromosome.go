package easyga

import (
	"math/rand"
)

type Chromosome struct {
	gene    []byte
	fitness float64
}

func getRandomGenotype(genotype uint8) byte {
	return byte(rand.Intn(int(genotype)))
}

func getRandomGeneIndex(c *Chromosome) int {
	return rand.Int() % len(c.gene)
}

func (c *Chromosome) Random(length int, genotype uint8) {
	c.gene = make([]byte, 0)
	for i := 0; i < length; i++ {
		c.gene = append(c.gene, getRandomGenotype(genotype))
	}
}

func (c *Chromosome) UpdateFitness() {
	// TODO: Make fitness function be customizable
	c.fitness = 0
	for _, genotype := range c.gene {
		c.fitness += float64(genotype)
	}
}

func (c *Chromosome) Mutate(genotype uint8) {
	c.gene[getRandomGeneIndex(c)] = getRandomGenotype(genotype)
	c.UpdateFitness()
}

func (c *Chromosome) Crossover(parent2 Chromosome) (child1, child2 Chromosome) {
	parent1 := c
	length := len(parent1.gene)
	position := getRandomGeneIndex(parent1)

	child1 = Chromosome{gene: make([]byte, length)}
	child2 = Chromosome{gene: make([]byte, length)}

	copy(child1.gene, parent1.gene[:position])
	copy(child2.gene, parent2.gene[:position])
	child1.gene = append(child1.gene[:position], parent2.gene[position:]...)
	child2.gene = append(child2.gene[:position], parent1.gene[position:]...)

	child1.UpdateFitness()
	child2.UpdateFitness()

	return child1, child2
}
