package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
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

	startPosition := 50

	realPass := 0
	switch *puzzlePart {
	case 1:
		realPass = PuzzlePart1(startPosition, file)
	case 2:
		realPass = PuzzlePart2(startPosition, file)
	default:
		fmt.Println("Invalid part number! Please enter 1 or 2")
	}

	fmt.Printf("Real password found to be %d\n", realPass)
}

func PuzzlePart1(startPosition int, file *os.File) int {
	realPass := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Printf("Current password %d\n", realPass)
		fmt.Printf("currentPosition %d\n", startPosition)
		line := scanner.Text()
		fmt.Println(line)
		direction := rune(line[0])
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			fmt.Printf("Error while parsing line! %v\n", err)
			continue
		}

		if direction == 'R' {
			startPosition += count
		} else {
			startPosition -= count
		}

		startPosition = ClampToRange(0, 100, startPosition%100)
		if startPosition == 0 {
			realPass++
		}
		fmt.Println()
	}

	if scanErr := scanner.Err(); scanErr != nil {
		fmt.Printf("Error while scanning file! %v\n", scanErr)
		return -1
	}
	return realPass
}

// Apparently there isn't a native way for this in go? I'm surprised
func absValInt(x int) int {
	if x > 0 {
		return x
	}
	return -1 * x
}

// Pass end position BEFORE modulus
func PassedZeroCount(startPosition int, endPosition int) int {
	direction := (endPosition - startPosition) / absValInt(endPosition-startPosition)
	count := 0
	// There is 100% a better way to do this but I cba right now it's midnight
	for x := startPosition; x != endPosition; x += direction {
		//fmt.Printf("x=%d\n", x)
		if x%100 == 0 {
			count++
		}
	}
	return count
}

func ClampToRange(startRange int, endRange int, value int) int {
	if value < startRange {
		diff := startRange - value
		return endRange - diff
	} else if value >= endRange {
		diff := value - endRange
		return startRange + diff
	}
	return value
}

func PuzzlePart2(startPosition int, file *os.File) int {
	realPass := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Printf("Current password %d\n", realPass)
		line := scanner.Text()
		fmt.Printf("currentPosition %d\n", startPosition)
		fmt.Println(line)
		direction := rune(line[0])
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			fmt.Printf("Error while parsing line! %v\n", err)
			continue
		}

		newPosition := startPosition

		if direction == 'R' {
			newPosition += count
		} else {
			newPosition -= count
		}

		realPass += PassedZeroCount(startPosition, newPosition)

		startPosition = ClampToRange(0, 100, newPosition%100)
		fmt.Println()
	}

	if scanErr := scanner.Err(); scanErr != nil {
		fmt.Printf("Error while scanning file! %v\n", scanErr)
		return -1
	}
	return realPass
}
