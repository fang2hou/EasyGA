package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

func main() {
	file, err := os.Open("lena.png")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get bounds
	rect := img.Bounds()
	// color reduction
	for i := 0; i < rect.Max.Y; i++ {
		for j := 0; j < rect.Max.X; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			fmt.Printf("(%d,", r)
			fmt.Printf("%d,", g)
			fmt.Printf("%d) ", b)
		}
		fmt.Println("")
	}
}
