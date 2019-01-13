package easyga

import "errors"

// GeneticAlgorithmParameters is a struct that contains everything generation need.
type GeneticAlgorithmParameters struct {
	CrossoverProbability float64
	MutationProbability  float64
	PopulationSize       int
	ChromosomeLength     int
	IterationsLimit      int
	GenotypeNumber       int
	RandomSeed           int64
	UseRoutine           bool
}

// check method will return error when the given parameters is risky.
func (param *GeneticAlgorithmParameters) check() error {
	if param.CrossoverProbability < 0 || param.CrossoverProbability > 1 {
		return errors.New("Error: CrossoverProbability should be in [0, 1]")
	}
	if param.MutationProbability < 0 || param.MutationProbability > 1 {
		return errors.New("Error: MutationProbability should be in [0, 1]")
	}
	if param.PopulationSize <= 2 {
		return errors.New("Error: PopulationSize should > 2")
	}
	if param.GenotypeNumber <= 1 {
		return errors.New("Error: The number of genotype should > 1")
	}
	if param.ChromosomeLength <= 0 {
		return errors.New("Error: ChromosomeLength should > 0")
	}
	if param.IterationsLimit <= 0 {
		return errors.New("Error: IterationsLimit should > 0")
	}
	return nil
}
