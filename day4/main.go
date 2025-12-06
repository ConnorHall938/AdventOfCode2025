package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

	numberAccessible := 0
	switch *puzzlePart {
	case 1:
		numberAccessible = Puzzle(file)
	case 2:
		numberAccessible = Puzzle(file)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Sum of available rolls found to be %d\n", numberAccessible)
}

var mapWidth int = -1  // This is assigned, -1 is fine so we know if it doesn't get assigned
var mapHeight int = -1 // This gets counted, and used for indexing. If we start at 0 then we "have a row" before we start reading lines

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

type MapPoint struct {
	x int
	y int
}
type CheckMap map[MapPoint]bool

// Which cells should be checked? Account for edges
func createCheckMap(x, y int) CheckMap {
	checkMap := make(map[MapPoint]bool)
	checkMap[MapPoint{x - 1, y - 1}] = true
	checkMap[MapPoint{x - 1, y}] = true
	checkMap[MapPoint{x - 1, y + 1}] = true
	checkMap[MapPoint{x, y - 1}] = true
	checkMap[MapPoint{x, y + 1}] = true
	checkMap[MapPoint{x + 1, y - 1}] = true
	checkMap[MapPoint{x + 1, y}] = true
	checkMap[MapPoint{x + 1, y + 1}] = true

	// Left side
	if x == 0 {
		delete(checkMap, MapPoint{x - 1, y - 1})
		delete(checkMap, MapPoint{x - 1, y})
		delete(checkMap, MapPoint{x - 1, y + 1})
	}
	// Top
	if y == 0 {
		delete(checkMap, MapPoint{x - 1, y - 1})
		delete(checkMap, MapPoint{x, y - 1})
		delete(checkMap, MapPoint{x + 1, y - 1})
	}
	// Right side
	if x == mapWidth {
		delete(checkMap, MapPoint{x + 1, y - 1})
		delete(checkMap, MapPoint{x + 1, y})
		delete(checkMap, MapPoint{x + 1, y + 1})
	}
	// Bottom
	if y == mapHeight {
		delete(checkMap, MapPoint{x - 1, y + 1})
		delete(checkMap, MapPoint{x, y + 1})
		delete(checkMap, MapPoint{x + 1, y + 1})
	}

	return checkMap
}

func countAvailableRolls(rollmap string) int {
	totalCount := 0
	cellsChecked := 0
	for x := 0; x <= mapWidth; x++ {
		for y := 0; y <= mapHeight; y++ {
			// For each cell, if it's a roll
			if rollAtPosition(rollmap, x, y) {
				// Which cells should we check? Account for edges
				checkMap := createCheckMap(x, y)
				// Check the cells around it
				rollCount := 0
				for coordinate := range checkMap {
					xCheck, yCheck := coordinate.x, coordinate.y
					if rollAtPosition(rollmap, xCheck, yCheck) {
						rollCount++
					}
				}
				if rollCount < 4 {
					totalCount += 1
				}
			}
			cellsChecked++

		}
	}
	fmt.Println(cellsChecked)
	return totalCount
}

func rollAtPosition(rollmap string, x, y int) bool {
	return rollmap[x+y*mapWidth] == '@'
}

func Puzzle(file *os.File) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	rollMap := ""
	// Input the lines
	for scanner.Scan() {
		line := scanner.Text()
		// For the love of god this should be constant.
		mapWidth = len(line) - 1
		rollMap += line
		mapHeight++
	}

	fmt.Printf("mapSize %d, %d\n", mapWidth, mapHeight)

	runningTotal = countAvailableRolls(rollMap)

	return runningTotal
}
