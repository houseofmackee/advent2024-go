package main

import (
	"fmt"
	"os"
	"sort"
)

// I am lazy...
var pl = fmt.Println

// function to get the absolute difference betwen two ints
func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func main() {

	// load the file called input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		pl("Error: ", err)
		return
	}
	defer file.Close()

	// read the file into two slices, the input is a list of integers with two
	// numbers per line, separated by spaces
	var leftNums, rightNums []int
	for {
		var a, b int
		_, err := fmt.Fscanf(file, "%d %d\n", &a, &b)
		if err != nil {
			break
		}
		leftNums = append(leftNums, a)
		rightNums = append(rightNums, b)
	}

	// sort both slices
	sort.Slice(leftNums, func(i, j int) bool {
		return leftNums[i] < leftNums[j]
	})
	sort.Slice(rightNums, func(i, j int) bool {
		return rightNums[i] < rightNums[j]
	})

	// caluclate the sum of the absolute differences between the two slices
	var sum int
	for i := 0; i < len(leftNums); i++ {
		sum += absDiffInt(leftNums[i], rightNums[i])
	}
	pl("Total distance:", sum)

	var sscore int
	for i := 0; i < len(leftNums); i++ {
		for j := 0; j < len(rightNums); j++ {
			if leftNums[i] == rightNums[j] {
				sscore += leftNums[i]
			}
		}
	}

	pl("Similarity score:", sscore)
}
