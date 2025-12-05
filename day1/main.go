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
	if *puzzlePart == 1 {
		realPass = PuzzlePart1(startPosition, file)
	} else if *puzzlePart == 2 {
		realPass = PuzzlePart2(startPosition, file)
	} else {
		fmt.Println("Invalid part number! Please enter 1 or 2")
	}

	fmt.Printf("Real password found to be %d\n", realPass)
}

func PuzzlePart1(startPosition int, file *os.File) int {
	currentPosition := startPosition
	realPass := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		fmt.Printf("%c\n", line[0])
		fmt.Printf("%v\n", line[1:])
		direction := rune(line[0])
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			fmt.Printf("Error while parsing line! %v\n", err)
			continue
		}

		if direction == 'R' {
			currentPosition += count
		} else {
			currentPosition -= count
		}

		currentPosition %= 100
		if currentPosition == 0 {
			realPass++
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		fmt.Printf("Error while scanning file! %v\n", scanErr)
		return -1
	}
	return realPass
}
