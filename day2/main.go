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

	sumOfBad := 0
	switch *puzzlePart {
	case 1:
		sumOfBad = PuzzlePart1(file)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of bad IDs found to be %d\n", sumOfBad)
}

// Return 1 or 0 for true or false, idk just want to
func IsBadID(id int) int {
	strID := strconv.Itoa(id)
	if len(strID)%2 == 0 { //If it's even length, we should check
		midpoint := (len(strID) / 2) // Maybe not -1
		// fmt.Printf("Midpoint = %d, First half = (%v), Second half = (%v)", midpoint, id[:midpoint], id[midpoint:])
		if strID[midpoint:] == strID[:midpoint] {
			fmt.Printf("%d found to be naughty ID!\n", id)
			return 1
		}
	}
	return 0
}

func CountBetweenRange(startRange, endRange int) int {
	count := 0
	for i := startRange; i <= endRange; i++ {
		count += i * IsBadID(i)
	}
	return count
}

func PuzzlePart1(file *os.File) int {
	scanner := bufio.NewScanner(file)
	// Get a list of first/last IDs
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.IndexByte(string(data), ','); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If no comma found and at EOF, return the rest of the data as a token
		if atEOF {
			return len(data), data, nil
		}
		// Request more data
		return 0, nil, nil
	})

	runningTotal := 0
	// For each one, split it by the '-'
	for scanner.Scan() {
		idRange := scanner.Text()
		// fmt.Printf("Processing line %v\n", line)
		// fmt.Println(line)
		IDs := strings.Split(idRange, "-")
		// fmt.Printf("Midpoint = %d, First half = (%v), Second half = (%v)", midpoint, id[:midpoint], id[midpoint:])
		startRange, err := strconv.Atoi(IDs[0])
		if err != nil {
			fmt.Printf("Error parsing %v: %v\n", IDs[0], err)
			continue
		}
		endRange, err := strconv.Atoi(IDs[1])
		if err != nil {
			fmt.Printf("Error parsing %v: %v\n", IDs[1], err)
			continue
		}
		runningTotal += CountBetweenRange(startRange, endRange)
	}

	return runningTotal
}
