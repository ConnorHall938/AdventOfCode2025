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

	sumOfJolt := 0
	switch *puzzlePart {
	case 1:
		sumOfJolt = Puzzle(file, getMaxJoltSize2)
	case 2:
		sumOfJolt = Puzzle(file, getMaxJoltSize12)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of jolts found to be %d\n", sumOfJolt)
}

// Terrible name
// for example, battery bank 818181911112111
// MaxJoltsForNAsStart(bank, n-1) would return 11, as the second last 1 is selected, therefore the last must also be selected
// MaxJoltsForNAsStart(bank, n-2) would also return 11, as the 3rd last 1 is selected, and the max of remaining digits is 1
// MaxJoltsForNAsStart(bank, n-10) would return 89m as the 11th last digit is selected (8) and the max of the remaining digits is 9

func MaxJoltsSize2(bank []int) int {
	firstDigit := bank[0]
	secondDigit := slices.Max(bank[1:])
	return firstDigit*10 + secondDigit
}

// Ahhh why does no language natively support this?!
func IntegerPow(num, exponent int) int {
	//IDC about 0  because we will not use that here.
	for i := range exponent - 1 {
		num *= num
		i = i
	}
	return num
}

func MaxJoltsForMaxSizeN(bank []int, startPoint, N int) int {
	bankSize := len(bank)
	latestSelected := bankSize - N // The latest we need to pick, e.g. if we have 4 elements [0,1,2,3] and want to choose 3, we need to pick all after index 1

	// If the start point does not allow us to select enough, return 0
	if startPoint > latestSelected {
		return 0
	}
	if N == 2 {
		overall := getMaxJoltSize2(bank[startPoint:])
		fmt.Printf("Max Jolts for %v with size %d (basecase) = %d\n", bank[startPoint:], N, overall)
		return overall
	}
	// If N == length, return all selected
	if N == len(bank[startPoint:]) {
		sum := 0
		for _, val := range bank[startPoint:] {
			sum += val
			sum *= 10
		}
		// We multiplied one too many times
		sum /= 10
		fmt.Printf("Max Jolts for %v with size %d = %d\n", bank[startPoint:], N, sum)
		return sum
	}

	maxSizeNMinus1 := bank[startPoint]*IntegerPow(10, N-1) + MaxJoltsForMaxSizeN(bank, startPoint+1, N-1)
	maxSizeIgnoreCurrent := MaxJoltsForMaxSizeN(bank, startPoint+1, N)
	maxOverall := IntMax(maxSizeNMinus1, maxSizeIgnoreCurrent)
	fmt.Printf("Max Jolts for %v with size %d = %d\n", bank[startPoint:], N, maxOverall)

	return maxOverall
}

func getMaxJoltSize12(bank []int) int {
	return MaxJoltsForMaxSizeN(bank, 0, 12)
}

func IntMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func getMaxJoltSize2(bank []int) int {

	currentMax := 0
	for i := len(bank) - 2; i >= 0; i-- {
		maxForSize := MaxJoltsSize2(bank[i:])
		fmt.Printf("MaxJolt for bank %v with %c as the first digit = %d\n", bank, bank[i], maxForSize)
		currentMax = IntMax(currentMax, maxForSize)
	}

	return currentMax
}

func Puzzle(file *os.File, joltFunc MaxJoltsForNAsStart) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	// For each one, split it by the '-'
	for scanner.Scan() {
		bankValue := scanner.Text()
		fmt.Printf("Processing batteries %v\n", bankValue)
		// Parse to list of ints for ease
		intList := make([]int, 0, len(bankValue))
		for _, rune := range bankValue {
			digit, err := strconv.Atoi(string(rune))
			if err != nil {
				fmt.Printf("Error converting character %c to int: %v\n", rune, err)
				return -1 // Handle the error appropriately
			}
			intList = append(intList, digit)
		}

		jolts := joltFunc(intList)
		runningTotal += jolts
		fmt.Printf("MaxJolts for bank %v = %d\n", bankValue, jolts)
	}

	return runningTotal
}
