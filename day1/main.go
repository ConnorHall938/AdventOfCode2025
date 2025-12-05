package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const puzzle_file = "puzzle.txt"

func main() {
	file, err := os.Open(puzzle_file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		fmt.Printf("Error opening puzzle file! Exiting!")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	currentPosition := 50
	realPass := 0

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
		fmt.Printf("Error while scanning file! %v\n", err)
		return
	}

	fmt.Printf("Real password found to be %d\n", realPass)
}
