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

const puzzleFileDefault = "./puzzle.txt"

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

	numberBoxes := 0
	switch *puzzlePart {
	case 1:
		numberBoxes = Puzzle(file)
	case 2:
		numberBoxes = Puzzle(file)
	default:
		fmt.Println("Invalid part number! Please enter 1")
		return
	}

	fmt.Printf("Total result is %d\n", numberBoxes)
}

// Ahhh why does no language natively support this?!
func IntegerPow(num, exponent int) int {
	curVal := 1
	for range exponent {
		curVal *= num
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

type Point struct {
	pointNumber, x, y int
}

type pointPair struct {
	p1, p2 *Point
}

type pointPairArea struct {
	area  int
	pairs *pointPair
}

func AreaBetweenPoints(p1, p2 Point) int {
	xDiff := p1.x - p2.x
	yDiff := p1.y - p2.y

	if p1.x < p2.x {
		xDiff = p2.x - p1.x
	}
	if p1.y < p2.y {
		yDiff = p2.y - p1.y
	}
	area := xDiff * yDiff
	// Area includes their x-y lines
	area += xDiff + yDiff + 1
	return area
}

func Puzzle(file *os.File) int {
	scanner := bufio.NewScanner(file)
	pointCount := 0

	pointList := []Point{}
	for scanner.Scan() {
		fmt.Printf("Scanning point %d\n", pointCount)
		line := scanner.Text()

		coordinates := strings.Split(line, ",")

		vals := make([]int, 2)
		for idx, coordinate := range coordinates {
			coordinateVal := 0
			var err error
			if coordinate[0] != ' ' {
				coordinateVal, err = strconv.Atoi(coordinate)
			} else {
				coordinateVal, err = strconv.Atoi(coordinate[1:])
			}
			if err != nil {
				fmt.Printf("Error converting string %v to int: %v\n", coordinateVal, err)
				return -1 // Handle the error appropriately
			}
			vals[idx] = coordinateVal
		}

		// Each group starts as an element on its own.
		pointList = append(pointList, Point{pointCount, vals[0], vals[1]})
		pointCount++
	}

	// Find area for each pair of points
	pointPairAreas := map[int]*pointPair{}
	for i := 0; i < len(pointList)-1; i++ {
		point1 := &pointList[i]
		for j := i + 1; j < len(pointList); j++ {
			point2 := &pointList[j]
			pairArea := AreaBetweenPoints(*point1, *point2)
			pointPairAreas[pairArea] = &pointPair{point1, point2}
		}
	}

	// Sort the areas
	pointPairAreasList := []pointPairArea{}
	for key, val := range pointPairAreas {
		pointPairAreasList = append(pointPairAreasList, pointPairArea{key, val})
	}
	slices.SortFunc(pointPairAreasList, func(a, b pointPairArea) int { return int(b.area - a.area) })

	return pointPairAreasList[0].area
}
