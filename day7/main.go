package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const puzzleFileDefault = "puzzle.txt"

func main() {

	puzzlePart := flag.Int("Part", 1, "Which part of the puzzle? 1/2")
	puzzleFile := flag.String("Input", puzzleFileDefault, "Name of the puzzle file")

	flag.Parse()

	if puzzleFile == nil {
		fmt.Printf("Puzzle file cannot be empty!")
		return
	}

	if puzzlePart == nil {
		fmt.Printf("Puzzle part cannot be empty!")
		return
	}

	file, err := os.Open(*puzzleFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		fmt.Printf("Error opening puzzle file! Exiting!")
		return
	}
	defer file.Close()

	numberFresh := 0
	switch *puzzlePart {
	case 1:
		numberFresh = Puzzle(file)
	case 2:
		numberFresh = Puzzle(file)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Total fresh ingredients found to be %d\n", numberFresh)
}

// Ahhh why does no language natively support this?!
func IntegerPow(num, exponent int) int {
	curVal := 1
	for i := range exponent {
		curVal *= num
		i = i
	}
	return curVal
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func IntMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Edits the input beam map, and returns the number of splitters hit
func beamsFromLine(currentBeams *map[int]struct{}, line string) int {
	beamsToRemove := []int{}
	beamsToAdd := []int{}
	splittersHit := 0

	// Hack for the first row
	for idx, val := range line {
		if val == 'S' {
			beamsToAdd = append(beamsToAdd, idx)
		}
	}

	for beamIndex := range *currentBeams {
		if line[beamIndex] == '^' {
			beamsToAdd = append(beamsToAdd, beamIndex-1, beamIndex+1)
			beamsToRemove = append(beamsToRemove, beamIndex)
			splittersHit++
		}
		//Beam source
		if line[beamIndex] == 'S' {
			beamsToAdd = append(beamsToAdd, beamIndex)
		}
	}

	// Remove then add
	// Remove
	for _, x := range beamsToRemove {
		delete(*currentBeams, x)
	}
	// Add
	for _, x := range beamsToAdd {
		(*currentBeams)[x] = struct{}{}
	}

	return splittersHit
}

func Puzzle(file *os.File) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	debugLineCounter := 0
	currentBeams := make(map[int]struct{}) // Column indexes for where currently active beams are
	for scanner.Scan() {
		fmt.Printf("Scanning line %d of beams\n", debugLineCounter)
		line := scanner.Text()
		runningTotal += beamsFromLine(&currentBeams, line)
		debugLineCounter++
	}

	return runningTotal
}
