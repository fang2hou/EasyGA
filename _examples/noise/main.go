package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"

	"github.com/fang2hou/easyga"
)

type imageData struct {
	yBound, xBound int
	brightness     [][]int
}

type imageNoise struct {
	ga       easyga.GeneticAlgorithm
	imageGap imageData
}

func main() {
	var findNoise imageNoise

	// Original
	findNoise.imageGap = getBrightnessFromFile("lena.png")
	// Corrupted
	corruptedImage := getBrightnessFromFile("lena.png_noisy_NA_60.0000_NFRow_0.0100_NFCol_0.0100.png")

	// Cache gap for better performance
	for i := 0; i < findNoise.imageGap.xBound; i++ {
		for j := 0; j < findNoise.imageGap.yBound; j++ {
			findNoise.imageGap.brightness[i][j] -= corruptedImage.brightness[i][j]
		}
	}

	fmt.Println(findNoise.totalError(60.0, 0.01, 0.01))

}

func (in *imageNoise) init() {
	// TODO: configuration
	parameters := easyga.GeneticAlgorithmParameters{
		CrossoverProbability: .8,
		MutationProbability:  .2,
		PopulationSize:       20,
		GenotypeNumber:       2,
		ChromosomeLength:     10,
		IterationsLimit:      1000,
		RandomSeed:           42,
	}

	custom := easyga.GeneticAlgorithmFunctions{
		ChromosomeInitFunction: func(c *easyga.Chromosome) {
			return
		},
		FitnessFunction: func(c *easyga.Chromosome) {
			return
		},
		CheckStopFunction: func(ga *easyga.GeneticAlgorithm) bool {
			return false
		},
	}

	if err := in.ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}
}

func (in *imageNoise) totalError(NoiseAmp float64, NoiseFreqRow float64, NoiseFreqCol float64) (fitness float64) {
	var noise float64

	// Error function assumed
	// N(row, col) = NoiseAmp×sin([2π×NoiseFreqRow×row]+[2π×NoiseFreqCol×col])
	for row := 0; row < in.imageGap.yBound; row++ {
		for col := 0; col < in.imageGap.xBound; col++ {
			noise = NoiseAmp * math.Sin(2.0*math.Pi*(NoiseFreqRow*float64(row+1)+NoiseFreqCol*float64(col+1)))
			fitness += float64(in.imageGap.brightness[row][col]) + noise
		}
	}

	return
}

func getBrightnessFromFile(fileName string) (inputImage imageData) {
	// Open image file
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode
	img, _, err := image.Decode(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Get bounds
	rect := img.Bounds()
	inputImage.xBound = rect.Max.X
	inputImage.yBound = rect.Max.Y

	// Color reduction
	for i := 0; i < rect.Max.Y; i++ {
		tempLine := make([]int, 0)
		for j := 0; j < rect.Max.X; j++ {
			// For example image, red = green = blue = brightness
			red, _, _, _ := img.At(j, i).RGBA()
			// Convert [0, 65536) -> [0, 256)
			tempLine = append(tempLine, int(red>>8))
		}
		inputImage.brightness = append(inputImage.brightness, tempLine)
	}

	return
}
