package main

import (
	"bufio"
	"fmt"
	"github.com/fang2hou/easyga"
	"image"
	"image/png"
	"math"
	"os"
)

func main() {
	var ga easyga.GeneticAlgorithm

	precision := 8 // equal pow2

	parameters := easyga.Parameters{
		CrossoverProbability: 0.8,
		MutationProbability:  .05,
		PopulationSize:       10,
		Genotype:             2,
		ChromosomeLength:     precision * 3,
		IterationsLimit:      10,
	}
	originalImage, err := decodePNG("lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	noisyImage, err := decodePNG("lena_noisy.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(originalImage)
	custom := easyga.CustomFunctions{}

	//custom.ChromosomeInitFunction = func(c *easyga.Chromosome) {
	//	You can customize your fitness function here
	//}

	//custom.MutateFunction = func(c *easyga.Chromosome) {
	//	You can customize your crossover function here
	//}

	custom.FitnessFunction = func(c *easyga.Chromosome) {
		c.Fitness = 0.0
		parameterBinaryLength := parameters.ChromosomeLength
		noiseAmp := make([]string, 0)
		noiseFreqRow := make([]string, 0)
		noiseFreqCol := make([]string, 0)

		//–  NoiseAmp 0 to 30.0
		//–  NoiseFreqRow 0 to 0.01
		//–  NoiseFreqCol 0 to 0.01
		for i := 0; i < parameterBinaryLength; i++ {
			if i < parameterBinaryLength/3 {
				if c.Gene[i] == 1 {
					noiseAmp = append(noiseAmp, "1")
				} else if c.Gene[i] == 0{
					noiseAmp = append(noiseAmp, "0")
				}
			} else if i >= parameterBinaryLength/3 && i < parameterBinaryLength/3*2 {
				if c.Gene[i] == 1 {
					noiseFreqRow = append(noiseFreqRow, "1")
				} else if c.Gene[i] == 0{
					noiseFreqRow = append(noiseFreqRow, "0")
				}
			} else if i >= parameterBinaryLength/3*2 && i < parameterBinaryLength {
				if c.Gene[i] == 1 {
					noiseFreqCol = append(noiseFreqCol, "1")
				} else if c.Gene[i] == 0{
					noiseFreqCol = append(noiseFreqCol, "0")
				}
			} else {
				fmt.Println("error")
				return
			}
		}

		c.Fitness = imageSimilarity(originalImage, noisyImage, noiseAmp, noiseFreqRow, noiseFreqCol, precision)
		if ga.Iteration%100 == 0 {
			outputImage(originalImage, noisyImage, ga.Iteration)
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

func decodePNG(filePath string) (image.Image, error) {
	reader, err := os.Open(filePath)
	fmt.Println(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer reader.Close()

	return png.Decode(bufio.NewReader(reader))
}

func encodePNG(img image.Image) (filePath string, err error) {
	filePath = "lena_output.png"
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return filePath, err
	}
	defer writer.Close()

	png.Encode(writer, img)
	fmt.Print(filePath)
	fmt.Print(img)
	return
}

//func addNoise(targetImage image.Image, noiseAmp []string, noiseFreqRow []string, noiseFreqCol []string, precision int) (resolvedImage image.Image) {
//	var NA, NFR, NFC float64
//	//–  NoiseAmp 0 to 30.0
//	//–  NoiseFreqRow 0 to 0.01
//	//–  NoiseFreqCol 0 to 0.01
//	for i := 0; i < precision; i++ {
//		if noiseAmp[i] == "1" {
//			NA += math.Pow(2.0, float64(i))
//		} else if noiseAmp[i] == "0" {
//		}
//
//		if noiseFreqRow[i] == "1" {
//			NFR += math.Pow(2.0, float64(i))
//		} else if noiseAmp[i] == "0" {
//		}
//
//		if noiseFreqCol[i] == "1" {
//			NFC += math.Pow(2.0, float64(i))
//		} else if noiseAmp[i] == "0" {
//		}
//	}
//	NA = 30.0 / math.Pow(2.0, float64(precision)) * NA
//	NFR = 30.0 / math.Pow(2.0, float64(precision)) * NFR
//	NFC = 30.0 / math.Pow(2.0, float64(precision)) * NFC
//
//	dx := targetImage.Bounds().Dx()
//	dy := targetImage.Bounds().Dy()
//
//	for i:=0;i<dx;i++{
//		for j := 0 ; j < dy;j++{
//			tempR,_,_,_ := targetImage.At(i,j).RGBA()
//			noise := NA*math.Sin(2.0*math.Pi*float64(dx)+2.0*math.Pi*float64(dy))
//			tempY := float64(tempR)/65535*256 + noise
//			targetImage.
//
//		}
//	}
//
//
//	return
//}

func imageSimilarity(targetImage image.Image, noisyImage image.Image, noiseAmp []string, noiseFreqRow []string, noiseFreqCol []string, precision int) (result float64) {
	var NA, NFR, NFC float64
	//–  NoiseAmp 0 to 30.0
	//–  NoiseFreqRow 0 to 0.01
	//–  NoiseFreqCol 0 to 0.01
	for i := 0; i < precision; i++ {
		if noiseAmp[i] == "1" {
			NA += math.Pow(2.0, float64(i))
		} else if noiseAmp[i] == "0"{
		}

		if noiseFreqRow[i] == "1" {
			NFR += math.Pow(2.0, float64(i))
		} else if noiseFreqRow[i] == "0"{
		}

		if noiseFreqCol[i] == "1" {
			NFC += math.Pow(2.0, float64(i))
		} else if noiseFreqCol[i] == "0"{
		}
	}
	NA = 30.0 / math.Pow(2.0, float64(precision)) * NA
	NFR = 0.01 / math.Pow(2.0, float64(precision)) * NFR
	NFC = 0.01 / math.Pow(2.0, float64(precision)) * NFC

	dx := targetImage.Bounds().Dx()
	dy := targetImage.Bounds().Dy()
	//resolvedImage := make()

	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			tempR, _, _, _ := targetImage.At(i, j).RGBA()
			noisyR, _, _, _ := noisyImage.At(i, j).RGBA()
			noise := NA * math.Sin(2.0*math.Pi*float64(dx)*NFR+2.0*math.Pi*float64(dy)*NFC)
			result += float64(tempR)/65535*255 + noise - float64(noisyR)/65535*255
			//append()
		}
	}

	return
}

func outputImage(originalImage image.Image, noisyImage image.Image, iteration int) {
	return
}
