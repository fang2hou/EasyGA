package easyga

// Chromosome is a struct that contains everything an individual need.
type Chromosome struct {
	Gene           []byte
	Fitness        float64
	GenotypeNumber int
}

// GetRandomGenotype method will return a random genotype.
func (c *Chromosome) GetRandomGenotype() byte {
	return byte(Rand.Intn(c.GenotypeNumber))
}

// GetRandomGeneIndex method will return a random index by the length of chromosome.
func (c *Chromosome) GetRandomGeneIndex() int {
	return Rand.Int() % len(c.Gene)
}

// Length method will return the length of chromosome.
func (c *Chromosome) Length() int {
	return len(c.Gene)
}

// Random method will generate individual randomly.
func (c *Chromosome) Random(length int, genotypeNumber int) {
	c.Gene = make([]byte, 0)
	for i := 0; i < length; i++ {
		c.Gene = append(c.Gene, c.GetRandomGenotype())
	}
}
