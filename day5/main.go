package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const puzzleFileDefault = "puzzle.txt"

type MaxJoltsForNAsStart func([]int) int

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

	fmt.Printf("Sum of fresh ingredients found to be %d\n", numberFresh)
}

// Ahhh why does no language natively support this?!
func IntegerPow(num, exponent int) int {
	//IDC about 0  because we will not use that here.
	curVal := num
	for i := range exponent - 1 {
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

func printSortedKeys(set map[int]bool) {
	keys := make([]int, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Println(keys)
}

func appendFreshIngredients(start, end int, list *map[int]bool) {
	for i := start; i <= end; i++ {
		(*list)[i] = true
	}
}

func Puzzle(file *os.File) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	freshIngredients := make(map[int]bool)
	// Input the lines
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		freshRange := strings.Split(line, "-")
		freshStart, err := strconv.Atoi(freshRange[0])
		if err != nil {
			fmt.Printf("Error converting string %v to int: %v\n", freshRange[0], err)
			return -1 // Handle the error appropriately
		}
		freshEnd, err := strconv.Atoi(freshRange[1])
		if err != nil {
			fmt.Printf("Error converting string %v to int: %v\n", freshRange[1], err)
			return -1 // Handle the error appropriately
		}
		appendFreshIngredients(freshStart, freshEnd, &freshIngredients)
	}

	printSortedKeys(freshIngredients)

	return runningTotal
}
