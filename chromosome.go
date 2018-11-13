package easyga

import (
	"math/rand"
)

// Chromosome is a struct that contains everything an individual need.
type Chromosome struct {
	Gene    []byte
	Fitness float64
}

func getRandomGenotype(genotype uint8) byte {
	return byte(rand.Intn(int(genotype)))
}

// Length method will return the length of chromosome
func (c *Chromosome) Length() int {
	return len(c.Gene)
}

// GetRandomGeneIndex method will return a random index by the length of chromosome
func (c *Chromosome) GetRandomGeneIndex() int {
	return rand.Int() % len(c.Gene)
}

// Random method will generate individual randomly.
func (c *Chromosome) Random(length int, genotype uint8) {
	c.Gene = make([]byte, 0)
	for i := 0; i < length; i++ {
		c.Gene = append(c.Gene, getRandomGenotype(genotype))
	}
}
