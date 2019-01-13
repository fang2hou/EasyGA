package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/fang2hou/easyga"
)

type imageData struct {
	height, width int
	brightness    [][]int
}

type imageNoise struct {
	ga                  easyga.GeneticAlgorithm
	original, corrupted imageData
	fitnessData         []float64
}

func main() {
	var findNoise imageNoise
	// Original
	findNoise.original = readImageFromFile("lena.png")

	// Corrupted
	findNoise.corrupted = readImageFromFile("noise.png")
	//findNoise.corrupted = readFromFile("noise_test.png") // Test: 60.0 0.01 0.01

	// Confirm the size of two pictures is same.
	if findNoise.corrupted.width != findNoise.original.width ||
		findNoise.corrupted.height != findNoise.original.width {
		fmt.Println("The size of two pictures is not same!")
		return
	}

	// Test
	// testNoiseAmp := 60.
	// testNoiseFreqRow := 0.01
	// testNoiseFreqCol := 0.01
	// fmt.Println(findNoise.totalError(testNoiseAmp, testNoiseFreqRow, testNoiseFreqCol))
	// return

	// Initialize a noise problem
	findNoise.init()

	// Run
	best, bestFit, iteration := findNoise.run()

	// Print out results
	NoiseAmp, NoiseFreqRow, NoiseFreqCol := genotypeToPhenotype(best.Gene)
	fmt.Println("Best gene is", NoiseAmp, NoiseFreqRow, NoiseFreqCol)
	fmt.Println("Best fitness is ", bestFit/(512*512))
	fmt.Println("Find it in", iteration, "generation.")

	drawFitnessChart(findNoise.fitnessData)
}

func (in *imageNoise) init() {
	// Every gene:
	// Index  0- 7 => NoiseAmp
	// Index  8-15 => NoiseFreqRow
	// Index 16-23 => NoiseFreqCol
	parameters := easyga.GeneticAlgorithmParameters{
		CrossoverProbability: .9,
		MutationProbability:  .1,
		PopulationSize:       200,
		GenotypeNumber:       2,
		ChromosomeLength:     24,
		IterationsLimit:      60,
		RandomSeed:           time.Now().UnixNano(),
		UseRoutine:           false,
	}

	custom := easyga.GeneticAlgorithmFunctions{
		FitnessFunction: func(c *easyga.Chromosome) {
			NoiseAmp, NoiseFreqRow, NoiseFreqCol := genotypeToPhenotype(c.Gene)
			// We want to find the individual with smallest fitness
			c.Fitness = -in.totalError(NoiseAmp, NoiseFreqRow, NoiseFreqCol)

			return
		},
		CheckStopFunction: func(ga *easyga.GeneticAlgorithm) bool {
			if ga.Population.Iteration >= ga.Parameters.IterationsLimit {
				return true
			}

			return false
		},
		StatisticFunction: func(ga *easyga.GeneticAlgorithm) {
			bestIndex, bestFitness := ga.Population.FindBest()
			in.fitnessData = append(in.fitnessData, bestFitness)
			if in.ga.Population.Iteration%10 == 0 && in.ga.Population.Iteration > 1 {
				NoiseAmp, NoiseFreqRow, NoiseFreqCol := genotypeToPhenotype(ga.Population.Chromosomes[bestIndex].Gene)
				in.outputOriginalImageWithNoiseCalculated(NoiseAmp, NoiseFreqRow, NoiseFreqCol)
			}
		},
	}

	if err := in.ga.Init(parameters, custom); err != nil {
		fmt.Println(err)
		return
	}
}

func (in *imageNoise) run() (easyga.Chromosome, float64, int) {
	return in.ga.Run()
}

// -----------------------------------------------------------------
// Gene translation Functions
// -----------------------------------------------------------------

func geneToPercent(gene []byte) (percent float64) {
	length := len(gene)
	max := math.Pow(2., float64(length)) - 1

	number := float64(0)
	for i := length - 1; i >= 0; i-- {
		number += float64(gene[i]) * math.Pow(2.0, float64(i))
	}

	percent = number / max

	return
}

func genotypeToPhenotype(gene []byte) (NoiseAmp float64, NoiseFreqRow float64, NoiseFreqCol float64) {
	NoiseAmp = geneToPercent(gene[:8]) * 30.
	NoiseFreqRow = geneToPercent(gene[8:16]) * .01
	NoiseFreqCol = geneToPercent(gene[16:24]) * .01

	return
}

func (in *imageNoise) totalError(NoiseAmp float64, NoiseFreqRow float64, NoiseFreqCol float64) (fitness float64) {
	// Error function assumed (Example)
	// NoiseAmp     (0, 30.0]
	// NoiseFreqRow (0, 0.01]
	// NoiseFreqCol (0, 0.01]
	// N(row, col) = NoiseAmp × sin([2π × NoiseFreqRow × row] + [2π × NoiseFreqCol × col])

	for row := 0; row < in.original.height; row++ {
		for col := 0; col < in.original.width; col++ {
			newBrightness := float64(in.original.brightness[row][col]) +
				NoiseAmp*math.Sin(2.0*math.Pi*(NoiseFreqRow*float64(row+1)+NoiseFreqCol*float64(col+1)))

			// Round
			newBrightness = math.Floor(newBrightness + 0.5)

			// Fix the new brightness over 255 or below 0.
			if newBrightness < 0 {
				newBrightness = 0.
			} else if newBrightness > 255 {
				newBrightness = 255.
			}

			// (Original + NoiseGA) - Corrupted
			fitness += math.Abs(newBrightness - float64(in.corrupted.brightness[row][col]))
		}
	}

	return
}

// -----------------------------------------------------------------
// I/O Functions
// -----------------------------------------------------------------

func readImageFromFile(fileName string) (inputImage imageData) {
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
	inputImage.width = rect.Max.X
	inputImage.height = rect.Max.Y

	// Color reduction
	for i := 0; i < rect.Max.Y; i++ {
		tempLine := make([]int, 0)
		for j := 0; j < rect.Max.X; j++ {
			// For grayscale image, red = green = blue = brightness.
			red, _, _, _ := img.At(j, i).RGBA()
			// Convert [0, 65536) -> [0, 256)
			tempLine = append(tempLine, int(red>>8))
		}
		inputImage.brightness = append(inputImage.brightness, tempLine)
	}

	return
}

func (in *imageNoise) outputImageWithNoiseCalculated(NoiseAmp float64, NoiseFreqRow float64, NoiseFreqCol float64) {
	filePath := "Noise_" + strconv.Itoa(in.ga.Population.Iteration) + ".png"

	// Open output image file
	outputFile, err := os.Create(filePath)
	defer outputFile.Close()

	if err != nil {
		fmt.Println("An error occurred when program try to read/create the output image!")
		return
	}

	// Create a new image
	outputImage := image.NewRGBA(image.Rect(0, 0, in.original.width, in.original.height))
	outputImageData := make([][]uint8, 0)

	// Calculation
	for row := 0; row < in.original.height; row++ {
		tempLine := make([]uint8, 0)
		for col := 0; col < in.original.width; col++ {
			newBrightness := float64(in.original.brightness[row][col]) +
				NoiseAmp*math.Sin(2.0*math.Pi*(NoiseFreqRow*float64(row+1)+NoiseFreqCol*float64(col+1)))
			// Round
			newBrightness = math.Floor(newBrightness + 0.5)

			// Fix the new brightness over 255 or below 0.
			if newBrightness < 0 {
				newBrightness = 0.
			} else if newBrightness > 255 {
				newBrightness = 255.
			}

			// Original + NoiseGA
			tempLine = append(tempLine, uint8(newBrightness))
		}
		outputImageData = append(outputImageData, tempLine)
	}

	// Convert grayscale to RGBA
	for row := 0; row < in.original.height; row++ {
		for col := 0; col < in.original.width; col++ {
			outputImage.SetRGBA(col, row, color.RGBA{outputImageData[row][col], outputImageData[row][col], outputImageData[row][col], 255})
		}
	}

	// Encode
	png.Encode(outputFile, outputImage)
}

func (in *imageNoise) outputOriginalImageWithNoiseCalculated(NoiseAmp float64, NoiseFreqRow float64, NoiseFreqCol float64) {
	filePath := "Original_" + strconv.Itoa(in.ga.Population.Iteration) + ".png"

	// Open output image file
	outputFile, err := os.Create(filePath)
	defer outputFile.Close()

	if err != nil {
		fmt.Println("An error occurred when program try to read/create the output image!")
		return
	}

	// Create a new image
	outputImage := image.NewRGBA(image.Rect(0, 0, in.corrupted.width, in.corrupted.height))
	outputImageData := make([][]uint8, 0)

	// Calculation
	for row := 0; row < in.corrupted.height; row++ {
		tempLine := make([]uint8, 0)
		for col := 0; col < in.corrupted.width; col++ {
			newBrightness := float64(in.corrupted.brightness[row][col]) -
				NoiseAmp*math.Sin(2.0*math.Pi*(NoiseFreqRow*float64(row+1)+NoiseFreqCol*float64(col+1)))
			// Round
			newBrightness = math.Floor(newBrightness + 0.5)

			// Fix the new brightness over 255 or below 0.
			if newBrightness < 0 {
				newBrightness = 0.
			} else if newBrightness > 255 {
				newBrightness = 255.
			}

			// Corrupted - NoiseGA
			tempLine = append(tempLine, uint8(newBrightness))
		}
		outputImageData = append(outputImageData, tempLine)
	}

	// Convert grayscale to RGBA
	for row := 0; row < in.corrupted.height; row++ {
		for col := 0; col < in.corrupted.width; col++ {
			outputImage.SetRGBA(col, row, color.RGBA{outputImageData[row][col], outputImageData[row][col], outputImageData[row][col], 255})
		}
	}

	// Encode
	png.Encode(outputFile, outputImage)
}
