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

	sumOfBad := 0
	switch *puzzlePart {
	case 1:
		sumOfBad = Puzzle(file, IsBadIDPart1)
	case 2:
		sumOfBad = Puzzle(file, IsBadIDPart2)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of bad IDs found to be %d\n", sumOfBad)
}

// Return 1 or 0 for true or false, idk just want to
func IsBadIDPart1(id int) int {
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

func isRepeating(str string, repeatLength int) bool {
	firstSegment := str[:repeatLength]
	// fmt.Printf("First segment of %v with repeatLength %d = %v\n", str, repeatLength, firstSegment)
	for i := repeatLength; i < len(str); i += repeatLength {
		// fmt.Printf("Comparing first segment %v with current segment %v\n", firstSegment, str[i:i+repeatLength])
		if str[i:i+repeatLength] != firstSegment {
			// fmt.Printf("%v is a good ID for repeatLength %d! has no repeating section of %v\n", str, repeatLength, firstSegment)
			return false
		}
	}
	fmt.Printf("%v is a bad ID! has a repeating section of %v\n", str, firstSegment)
	return true
}

func IsBadIDPart2(id int) int {
	strID := strconv.Itoa(id)
	for repeatLength := range len(strID) / 2 {
		if len(strID)%(repeatLength+1) == 0 && isRepeating(strID, repeatLength+1) {
			return 1
		}
		// fmt.Println()
	}
	return 0
}

func CountBetweenRange(startRange, endRange int, counterFunc idCounter) int {
	count := 0
	for i := startRange; i <= endRange; i++ {
		count += i * counterFunc(i)
	}
	return count
}

func Puzzle(file *os.File, counterFunc idCounter) int {
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
		runningTotal += CountBetweenRange(startRange, endRange, counterFunc)
	}

	return runningTotal
}
