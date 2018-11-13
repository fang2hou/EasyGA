package main

import (
	"bufio"
	"github.com/fang2hou/easyga"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	"image/png"
	"math"
	"os"
)


func main() {
	var ga easyga.GeneticAlgorithm

	imagick.Initialize()
	defer imagick.Terminate()
	precision := 8 // equal pow2

	parameters := easyga.Parameters{
		CrossoverProbability: 1,
		MutationProbability:  .1,
		PopulationSize:       4,
		Genotype:             2,
		ChromosomeLength:     precision * 3,
		IterationsLimit:      1000,
	}
	originalImage,err := decodePNG("lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	noisyImage,err := decodePNG("lena_noisy.png")
	if err != nil {
		fmt.Println(err)
		return
	}

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

		tempNoisyImage := addNoise(originalImage,noiseAmp,noiseFreqRow,noiseFreqCol,precision)
		c.Fitness = imageSimilarity(tempNoisyImage,noisyImage)
		if ga.Iteration % 100 == 0{
			outputImage(originalImage,noisyImage,ga.Iteration)
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

func addNoise(targetImage image.Image ,noiseAmp []string,noiseFreqRow []string,noiseFreqCol []string,precision int)(resolvedImage image.Image){
	var NA,NFR,NFC float64
	//–  NoiseAmp 0 to 30.0
	//–  NoiseFreqRow 0 to 0.01
	//–  NoiseFreqCol 0 to 0.01
	for i:=0;i<precision;i++{
		if noiseAmp[i] == "1"{
			NA += math.Pow(2.0,float64(i))
		} else if noiseAmp[i] == "0"{
		}

		if noiseFreqRow[i] == "1"{
			NFR += math.Pow(2.0,float64(i))
		} else if noiseAmp[i] == "0"{
		}

		if noiseFreqCol[i] == "1"{
			NFC += math.Pow(2.0,float64(i))
		} else if noiseAmp[i] == "0"{
		}
	}
	NA = 30.0 / math.Pow(2.0,float64(precision)) * NA
	NFR = 30.0 / math.Pow(2.0,float64(precision)) * NFR
	NFC = 30.0 / math.Pow(2.0,float64(precision)) * NFC
	
	return
}

func imageSimilarity(firstImage image.Image,secondImage image.Image)(result float64){
	return
}

func outputImage(originalImage image.Image,noisyImage image.Image,iteration int){
	return
}