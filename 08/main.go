package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	imgWidth := 25
	imgHeight := 6

	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	imgChars := scanner.Text()

	layerChunkSize := imgWidth * imgHeight

	maxNumZeroes := -1
	outputValueForBestLayer := -1

	var allLayers []string

	for len(imgChars) > 0 {
		layerData := imgChars[:layerChunkSize]
		imgChars = imgChars[layerChunkSize:]

		allLayers = append(allLayers, layerData)

		numZeroes := 0
		numOnes := 0
		numTwos := 0

		for _, c := range layerData {
			switch c {
			case '0':
				numZeroes += 1
			case '1':
				numOnes += 1
			case '2':
				numTwos += 1
			}
		}
		if maxNumZeroes < 0 || numZeroes < maxNumZeroes {
			maxNumZeroes = numZeroes
			outputValueForBestLayer = numOnes * numTwos
		}
	}
	fmt.Println("Part 1 solution: ", outputValueForBestLayer)

	/////////////////////////////////////////////////////////////
	// part 2: generate and output the image

	// Start with an image that is entirely transparent,
	// then apply each layer.
	image := strings.Repeat("2", layerChunkSize)

	for _, layer := range allLayers {
		newImage := ""
		for i, c := range image {
			if c == '2' {
				newImage += string(layer[i])
			} else {
				newImage += string(image[i])
			}
		}
		image = newImage
	}

	fmt.Println("Part 2 solution:\n")

	for len(image) > 0 {
		fmt.Println(strings.Replace(image[:imgWidth], "0", " ", -1))
		image = image[imgWidth:]
	}
}
