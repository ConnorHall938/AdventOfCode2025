package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
)

const puzzleFileDefault = "puzzle.txt"

type idCounter func(int) int

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

	sumOfJolt := 0
	switch *puzzlePart {
	case 1:
		sumOfJolt = Puzzle(file)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of jolts found to be %d\n", sumOfJolt)
}

// Terrible name
// for example, battery bank 818181911112111
// getMaxJoltWithNAsFirst(bank, n-1) would return 11, as the second last 1 is selected, therefore the last must also be selected
// getMaxJoltWithNAsFirst(bank, n-2) would also return 11, as the 3rd last 1 is selected, and the max of remaining digits is 1
// getMaxJoltWithNAsFirst(bank, n-10) would return 89m as the 11th last digit is selected (8) and the max of the remaining digits is 9

func getMaxJoltWithNAsFirst(bank []int, startPoint int) int {
	firstDigit := bank[startPoint]
	secondDigit := slices.Max(bank[startPoint+1:])
	return firstDigit*10 + secondDigit
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func getMaxJolt(bank string) int {
	// Parse to list of ints for ease
	intList := make([]int, 0, len(bank))
	for _, rune := range bank {
		digit, err := strconv.Atoi(string(rune))
		if err != nil {
			fmt.Printf("Error converting character %c to int: %v\n", rune, err)
			return -1 // Handle the error appropriately
		}
		intList = append(intList, digit)
	}

	currentMax := 0
	for i := len(bank) - 2; i >= 0; i-- {
		maxForSize := getMaxJoltWithNAsFirst(intList, i)
		fmt.Printf("MaxJolt for bank %v with %c as the first digit = %d\n", bank, bank[i], maxForSize)
		currentMax = IntMax(currentMax, maxForSize)
	}

	return currentMax
}

func Puzzle(file *os.File) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	// For each one, split it by the '-'
	for scanner.Scan() {
		bankValue := scanner.Text()
		fmt.Printf("Processing batteries %v\n", bankValue)

		jolts := getMaxJolt(bankValue)
		runningTotal += jolts
		fmt.Printf("MaxJolts for bank %v = %d\n", bankValue, jolts)
	}

	return runningTotal
}
