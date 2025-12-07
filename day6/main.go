package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		numberFresh = Puzzle(file, true)
	case 2:
		numberFresh = Puzzle(file, false)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Homework answer is %d\n", numberFresh)
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

func DoItWrong(file *os.File) int {

	return 0
}

func DoItNormally(file *os.File) int {
	scanner := bufio.NewScanner(file)
	lineCounter := 0
	line := "" // Litterally only outside the loop so I can break and still use the line
	numbers := [][]string{}
	// Parse the lines
	for scanner.Scan() {
		fmt.Printf("Scanning line %d of input\n", lineCounter)
		line = scanner.Text()
		if line[0] == '*' || line[0] == '+' {
			break
		}
		numberStrings := strings.Split(line, " ")
		filteredNumberStrings := []string{}
		for _, val := range numberStrings {
			if val != "" {
				filteredNumberStrings = append(filteredNumberStrings, val)
			}
		}

		numbers = append(numbers, make([]string, len(filteredNumberStrings)))
		copy(numbers[lineCounter], filteredNumberStrings)
		lineCounter++
	}

	// Parse the operation and calculate
	operationStrings := strings.Split(line, " ")
	filteredOperationStrings := []string{}
	for _, val := range operationStrings {
		if val != "" {
			filteredOperationStrings = append(filteredOperationStrings, val)
		}
	}

	return withHumanMath(numbers, filteredOperationStrings)
}

func withHumanMath(inputs [][]string, operands []string) int {
	runningTotal := 0
	for idx, val := range operands {
		currentTotal := 0
		switch val {
		case "*":
			currentTotal = 1
			for _, numbersList := range inputs {
				num, err := strconv.Atoi(numbersList[idx])
				if err != nil {
					fmt.Printf("Error converting string %v to int: %v\n", numbersList[idx], err)
					return -1 // Handle the error appropriately
				}
				currentTotal *= num
			}
		case "+":
			currentTotal = 0
			for _, numbersList := range inputs {
				num, err := strconv.Atoi(numbersList[idx])
				if err != nil {
					fmt.Printf("Error converting string %v to int: %v\n", numbersList[idx], err)
					return -1 // Handle the error appropriately
				}
				currentTotal += num
			}
		}
		runningTotal += currentTotal
	}
	return runningTotal
}

func Puzzle(file *os.File, fixTheirStupidAssMath bool) int {

	if fixTheirStupidAssMath {
		return DoItNormally(file)
	}
	return DoItWrong(file)

}
