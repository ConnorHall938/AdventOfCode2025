package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
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

	numberToConnect := 0
	switch *puzzleFile {
	case "./puzzle.txt":
		numberToConnect = 1000
	default:
		numberToConnect = 10
	}

	numberBoxes := 0
	switch *puzzlePart {
	case 1:
		numberBoxes = Puzzle(file, numberToConnect)
	case 2:
		numberBoxes = Puzzle(file, numberToConnect)
	default:
		fmt.Println("Invalid part number! Please enter 1")
		return
	}

	fmt.Printf("Total result is %d\n", numberBoxes)
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

type Box struct {
	boxNumber, groupNumber, x, y, z int
}

type boxPair struct {
	b1, b2 *Box
}

type boxPairDistance struct {
	distance float64
	boxes    *boxPair
}

func DistanceBetweenBoxes(box1, box2 Box) float64 {
	sumSquares := (box1.x-box2.x)*(box1.x-box2.x) + (box1.y-box2.y)*(box1.y-box2.y) + (box1.z-box2.z)*(box1.z-box2.z)
	return math.Sqrt(float64(sumSquares))
}

func Puzzle(file *os.File, numberBoxes int) int {
	scanner := bufio.NewScanner(file)
	boxCount := 0

	// List of which group each box belongs to
	boxList := []Box{}
	for scanner.Scan() {
		fmt.Printf("Scanning box %d\n", boxCount)
		line := scanner.Text()

		coordinates := strings.Split(line, ",")

		vals := make([]int, 3)
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
		boxList = append(boxList, Box{boxCount, boxCount, vals[0], vals[1], vals[2]})
		boxCount++
	}

	// Find distance for each pair of boxes
	boxPairDistances := map[float64]*boxPair{}
	for i := 0; i < len(boxList)-1; i++ {
		box1 := &boxList[i]
		for j := i + 1; j < len(boxList); j++ {
			box2 := &boxList[j]
			pairDistance := DistanceBetweenBoxes(*box1, *box2)
			boxPairDistances[pairDistance] = &boxPair{box1, box2}
		}
	}

	// Sort the distances
	boxPairDistancesList := []boxPairDistance{}
	for key, val := range boxPairDistances {
		boxPairDistancesList = append(boxPairDistancesList, boxPairDistance{key, val})
	}
	slices.SortFunc(boxPairDistancesList, func(a, b boxPairDistance) int { return int(a.distance - b.distance) })

	// Group them by closest
	groupSizes := make([]int, len(boxList))
	for i := range groupSizes {
		groupSizes[i] = 1
	}
	pairsGrouped := 0
	for i := range boxPairDistancesList {
		box1 := boxPairDistancesList[i].boxes.b1
		box2 := boxPairDistancesList[i].boxes.b2

		if box1.groupNumber == box2.groupNumber {
			// Same group - Nothing changes
		} else if groupSizes[box1.groupNumber] != 1 && groupSizes[box2.groupNumber] != 1 {
			// Merge them into group 1
			box2Group := box2.groupNumber
			for idx := range boxList {
				if boxList[idx].groupNumber == box2Group {
					boxList[idx].groupNumber = box1.groupNumber
					groupSizes[box1.groupNumber]++
				}
			}
			groupSizes[box2Group] = 0
		} else if groupSizes[box2.groupNumber] == 1 {
			groupSizes[box2.groupNumber] = 0
			(*box2).groupNumber = box1.groupNumber
			groupSizes[box1.groupNumber]++
		} else if groupSizes[box1.groupNumber] == 1 {
			groupSizes[box1.groupNumber] = 0
			(*box1).groupNumber = box2.groupNumber
			groupSizes[box2.groupNumber]++
		}

		pairsGrouped++
		if pairsGrouped >= numberBoxes {
			break
		}
	}

	slices.SortFunc(groupSizes, func(a, b int) int {
		return b - a
	})

	return groupSizes[0] * groupSizes[1] * groupSizes[2]

}
