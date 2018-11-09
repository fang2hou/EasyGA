package main

import (
	"bufio"
	"github.com/fang2hou/easyga"
	"fmt"
	"image"
	"image/png"
	"os"
)


func main() {
	var ga easyga.GeneticAlgorithm


	parameters := easyga.Parameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     24,
		IterationsLimit:      1000,
	}
	originalImage,err := decodePNG("lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	//noisyImage,err := decodePNG("lena_noisy.png")
	//if err != nil {
	//	return
	//}
	fmt.Println(originalImage)
	encodePNG(originalImage)

	custom := easyga.CustomFunctions{}

	//custom.ChromosomeInitFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.MutateFunction = func(c *easyga.Chromosome) {
	//	You can customize your crossover function here
	//}

	custom.FitnessFunction = func(c *easyga.Chromosome) {
		c.Fitness = 0
		parameterBinaryLength := parameters.ChromosomeLength
		noiseAmp := make([]string,0)
		noiseFreqRow := make([]string,0)
		noiseFreqCol := make([]string,0)

		//–  NoiseAmp 0 to 30.0
		//–  NoiseFreqRow 0 to 0.01
		//–  NoiseFreqCol 0 to 0.01
		for i := 0;i < parameterBinaryLength;i++ {
			if i < parameterBinaryLength / 3 {
				noiseAmp = append(noiseAmp,string(c.Gene[i]))
			} else if i >= parameterBinaryLength / 3  &&i < parameterBinaryLength / 3 * 2{
				noiseFreqRow = append(noiseFreqRow,string(c.Gene[i]))
			} else if i >= parameterBinaryLength / 3 * 2  &&i < parameterBinaryLength{
				noiseFreqCol = append(noiseFreqCol,string(c.Gene[i]))
			} else {
				fmt.Println("error")
				return
			}
		}

	}

	//custom.CrossOverFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.CheckStopFunction = func (ga *easyga.GeneticAlgorithm) bool {
	//
	//}

	if err := ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}

	best, bestFit, iteration := ga.Run()

	fmt.Println("Best gene is", best)
	fmt.Println("Best fitness is", bestFit)
	fmt.Println("Find it in", iteration, "generation.")
}

func decodePNG(filePath string)(image.Image,error) {
	reader, err := os.Open(filePath)
	fmt.Println(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer reader.Close()

	return png.Decode(bufio.NewReader(reader))
}

func encodePNG(img image.Image)(filePath string,err error) {
	filePath = "lena_output.png"
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return filePath,err
	}
	defer writer.Close()

	png.Encode(writer,img)
	fmt.Print(filePath)
	fmt.Print(img)
	return
}
