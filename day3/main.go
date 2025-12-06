package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
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

	sumOfJolt := 0
	switch *puzzlePart {
	case 1:
		sumOfJolt = Puzzle(file, 2)
	case 2:
		sumOfJolt = Puzzle(file, 12)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of jolts found to be %d\n", sumOfJolt)
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

// Just hack this in, cba restructuring again
type cacheKey struct {
	bank string
	n    int
}

var maxJoltsCache = make(map[cacheKey]int)

func arrayToString(bank []int) string {
	var builder strings.Builder
	for _, digit := range bank {
		builder.WriteString(strconv.Itoa(digit))
	}
	result := builder.String()
	return result
}

func MaxJoltsForMaxSizeN(bank []int, N int) int {
	key := cacheKey{arrayToString(bank), N}
	if val, ok := maxJoltsCache[key]; ok {
		return val
	}
	bankSize := len(bank)
	latestSelected := bankSize - N // The latest we need to pick, e.g. if we have 4 elements [0,1,2,3] and want to choose 3, we need to pick all after index 1

	// If the bank does not allow us to select enough, return 0
	if len(bank) < latestSelected {
		return 0
	}
	if N == 2 {
		overall := getMaxJoltSize2(bank)
		// fmt.Printf("Max Jolts for %v with size %d (basecase) = %d\n", bank[startPoint:], N, overall)
		maxJoltsCache[key] = overall
		return overall
	}
	// If N == length, return all selected
	if N == len(bank) {
		sum := 0
		for _, val := range bank {
			sum += val
			sum *= 10
		}
		// We multiplied one too many times
		sum /= 10
		// fmt.Printf("Max Jolts for %v with size %d = %d\n", bank[startPoint:], N, sum)
		maxJoltsCache[key] = sum
		return sum
	}

	maxSizeNMinus1 := bank[0]*IntegerPow(10, N-1) + MaxJoltsForMaxSizeN(bank[1:], N-1)
	maxSizeIgnoreCurrent := MaxJoltsForMaxSizeN(bank[1:], N)
	maxOverall := IntMax(maxSizeNMinus1, maxSizeIgnoreCurrent)
	maxJoltsCache[key] = maxOverall
	// fmt.Printf("Max Jolts for %v with size %d = %d\n", bank[startPoint:], N, maxOverall)

	return maxOverall
}

func MaxJoltsSize2(bank []int) int {
	firstDigit := bank[0]
	secondDigit := slices.Max(bank[1:])
	return firstDigit*10 + secondDigit
}

// Base case
func getMaxJoltSize2(bank []int) int {

	currentMax := 0
	for i := len(bank) - 2; i >= 0; i-- {
		maxForSize := MaxJoltsSize2(bank[i:])
		currentMax = IntMax(currentMax, maxForSize)
	}

	return currentMax
}

func getMaxJoltSizeN(bank []int, N int) int {
	maxJoltsForBank := MaxJoltsForMaxSizeN(bank, N)
	fmt.Printf("Max jolts for bank %v = %d\n", bank, maxJoltsForBank)
	return maxJoltsForBank
}

func Puzzle(file *os.File, numberAllowed int) int {
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
		jolts := getMaxJoltSizeN(intList, numberAllowed)
		runningTotal += jolts
		fmt.Printf("MaxJolts for bank %v = %d\n", bankValue, jolts)
	}

	return runningTotal
}
