package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
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
		numberFresh = Puzzle(file, false)
	case 2:
		numberFresh = Puzzle(file, true)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Total result is %d\n", numberFresh)
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

// Edits the input beam map, and returns the number of splitters hit
func beamsFromLine(currentBeams *map[int]struct{}, line string) int {
	beamsToRemove := []int{}
	beamsToAdd := []int{}
	splittersHit := 0

	// Hack for the first row
	for idx, val := range line {
		if val == 'S' {
			beamsToAdd = append(beamsToAdd, idx)
		}
	}

	for beamIndex := range *currentBeams {
		if line[beamIndex] == '^' {
			beamsToAdd = append(beamsToAdd, beamIndex-1, beamIndex+1)
			beamsToRemove = append(beamsToRemove, beamIndex)
			splittersHit++
		}
		//Beam source
		if line[beamIndex] == 'S' {
			beamsToAdd = append(beamsToAdd, beamIndex)
		}
	}

	// Remove then add
	// Remove
	for _, x := range beamsToRemove {
		delete(*currentBeams, x)
	}
	// Add
	for _, x := range beamsToAdd {
		(*currentBeams)[x] = struct{}{}
	}

	return splittersHit
}

func DoesBeamSplit(column int, line string) bool {
	return line[column] == '^'
}

// Edits the input beam worlds, potentially adding some.
func beamsFromLineManyWorlds(timelineColumns *[]int, line string) {
	newTimelineColumns := make([]int, len(*timelineColumns))
	copy(newTimelineColumns, *timelineColumns)
	// For each column
	for column, beamsInColumn := range *timelineColumns {
		if beamsInColumn == 0 {
			continue
		}

		beamSplits := DoesBeamSplit(column, line)
		if beamSplits {
			newTimelineColumns[column-1] += beamsInColumn
			newTimelineColumns[column+1] += beamsInColumn
			newTimelineColumns[column] = 0
		}
	}
	*timelineColumns = newTimelineColumns
}

func Puzzle(file *os.File, manyWorlds bool) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	debugLineCounter := 0
	if !manyWorlds {
		currentBeams := make(map[int]struct{}) // Column indexes for where currently active beams are
		for scanner.Scan() {
			fmt.Printf("Scanning line %d of beams\n", debugLineCounter)
			line := scanner.Text()
			runningTotal += beamsFromLine(&currentBeams, line)
			debugLineCounter++
		}
	} else {

		// Hack for the first row
		startIndex := 0
		lineLength := 0
		if scanner.Scan() {
			line := scanner.Text()
			lineLength = len(line)
			for idx, val := range line {
				if val == 'S' {
					startIndex = idx
					break
				}
			}
		} else {
			return -1
		}
		currentBeamWorlds := make([]int, lineLength) // Number of timelines in each column
		currentBeamWorlds[startIndex] = 1
		totalTimelines := 0
		for scanner.Scan() {
			lineStart := time.Now()

			fmt.Printf("Scanning line %d of beams\n", debugLineCounter)
			line := scanner.Text()
			beamsFromLineManyWorlds(&currentBeamWorlds, line)
			elapsed := time.Since(lineStart)

			totalTimelines = 0
			for _, val := range currentBeamWorlds {
				totalTimelines += val
			}

			fmt.Printf(
				"Line %d took %s — worlds: %d\n",
				debugLineCounter,
				elapsed,
				totalTimelines,
			)
			debugLineCounter++
		}
		runningTotal = totalTimelines
	}

	return runningTotal
}
