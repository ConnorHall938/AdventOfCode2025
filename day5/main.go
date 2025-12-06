package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

	numberFresh := 0
	switch *puzzlePart {
	case 1:
		numberFresh = Puzzle(file, true)
	case 2:
		numberFresh = Puzzle(file, false)
	default:
		fmt.Println("Invalid part number! Please enter 1")
	}

	fmt.Printf("Total fresh ingredients found to be %d\n", numberFresh)
}

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

func IntMin(x, y int) int {
	if x > y {
		return y
	}
	return x
}

type freshnessRange struct {
	start int
	end   int
}

// This one is AI because yeah... meaningless
func printSortedKeys(sets []freshnessRange) {
	// Make a copy so sorting does not mutate the original slice unless you want it to
	ranges := make([]freshnessRange, len(sets))
	copy(ranges, sets)

	// Sort by start, then by end
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].start == ranges[j].start {
			return ranges[i].end < ranges[j].end
		}
		return ranges[i].start < ranges[j].start
	})

	// Print the ranges
	for _, r := range ranges {
		fmt.Printf("[%d, %d]\n", r.start, r.end)
	}
}

func WithinRange(existingRange *freshnessRange, val int) bool {
	return val >= existingRange.start && val <= existingRange.end
}

func RemoveRangeFromList(list *[]freshnessRange, rangeToRemove freshnessRange) {
	for i, focusRange := range *list {
		if focusRange.start == rangeToRemove.start &&
			focusRange.end == rangeToRemove.end {
			*list = append((*list)[:i], (*list)[i+1:]...)
			return
		}
	}
}

func appendFreshIngredients(start, end int, list *[]freshnessRange) {
	// Gather ranges we start, end or encompass
	involvedRanges := []freshnessRange{}
	for i := range *list {
		existingRange := &(*list)[i]
		// Encompass the range
		if start < existingRange.start && end > existingRange.end {
			involvedRanges = append(involvedRanges, *existingRange)
		}
		// Start or end inside the range
		if WithinRange(existingRange, start) || WithinRange(existingRange, end) {
			involvedRanges = append(involvedRanges, *existingRange)
		}

	}

	if len(involvedRanges) == 0 {
		// This is a new range!
		*list = append(*list, freshnessRange{start, end})
		return
	}

	// This is a guess, BUT feels right
	// If we involve more than 1 range, we can just take the max of all the ends and mins of all the starts and merge into 1
	newRangeStart := start
	newRangeEnd := end
	for _, theInvolvedRange := range involvedRanges {
		newRangeStart = IntMin(newRangeStart, theInvolvedRange.start)
		newRangeEnd = IntMax(newRangeEnd, theInvolvedRange.end)
	}
	// Remove now we're done
	for _, theInvolvedRange := range involvedRanges {
		RemoveRangeFromList(list, theInvolvedRange)
	}
	*list = append(*list, freshnessRange{newRangeStart, newRangeEnd})
}

func IsIngredientFresh(ingredient int, freshIngredientsRanges []freshnessRange) bool {
	for _, freshRange := range freshIngredientsRanges {
		if ingredient >= freshRange.start && ingredient <= freshRange.end {
			return true
		}
	}
	return false
}

func countTotalFresh(freshIngredientsRanges []freshnessRange) int {
	total := 0

	for _, freshRange := range freshIngredientsRanges {
		total += freshRange.end - freshRange.start
		total += 1 // The range is inclusive, but the above cuts off an end
	}

	return total
}

func Puzzle(file *os.File, countPresentFresh bool) int {
	scanner := bufio.NewScanner(file)
	runningTotal := 0
	debugLineCounter := 0
	freshIngredients := []freshnessRange{}
	// Input the non-spoiled ingedients
	for scanner.Scan() {
		//fmt.Printf("Scanning line %d of freshness\n", debugLineCounter)
		line := scanner.Text()
		if line == "" {
			debugLineCounter = 0
			break
		}
		freshRange := strings.Split(line, "-")
		freshStart, err := strconv.Atoi(freshRange[0])
		if err != nil {
			fmt.Printf("Error converting string %v to int: %v\n", freshRange[0], err)
			return -1 // Handle the error appropriately
		}
		freshEnd, err := strconv.Atoi(freshRange[1])
		if err != nil {
			fmt.Printf("Error converting string %v to int: %v\n", freshRange[1], err)
			return -1 // Handle the error appropriately
		}
		appendFreshIngredients(freshStart, freshEnd, &freshIngredients)
		debugLineCounter++
	}

	printSortedKeys(freshIngredients)

	if countPresentFresh {
		for scanner.Scan() {
			fmt.Printf("Scanning line %d of ingredients - ", debugLineCounter)
			line := scanner.Text()
			ingredient, err := strconv.Atoi(line)
			if err != nil {
				fmt.Printf("Error converting string %v to int: %v\n", line, err)
				return -1 // Handle the error appropriately
			}

			if IsIngredientFresh(ingredient, freshIngredients) {
				runningTotal++
				fmt.Println("Fresh")
			} else {
				fmt.Println("Not fresh")
			}

			debugLineCounter++
		}
	} else {
		return countTotalFresh(freshIngredients)
	}

	return runningTotal
}
