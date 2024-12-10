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

	// open the file called input.txt
	file, err := os.Open("input_test.txt")
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

	// count unique numbers in the left slice and print the count
	var uniqueNums = make(map[int]bool)
	for _, leftNum := range leftNums {
		uniqueNums[leftNum] = true
	}
	pl("Unique numbers in the left slice:", len(uniqueNums))

	// count unique numbers in the right slice and print the count
	uniqueNums = make(map[int]bool)
	for _, rightNum := range rightNums {
		uniqueNums[rightNum] = true
	}
	pl("Unique numbers in the right slice:", len(uniqueNums))

	// calculate the sum of the absolute differences between the two slices
	var sum int
	for i := 0; i < len(leftNums); i++ {
		sum += absDiffInt(leftNums[i], rightNums[i])
	}
	pl("Total distance:", sum)

	// calculate the similarity score between the two slices
	var sscore int
	for _, leftNum := range leftNums {
		for _, rightNum := range rightNums {
			if leftNum == rightNum {
				sscore += leftNum
			}
		}
	}
	pl("Similarity score:", sscore)

	sscore = 0
	var freqMap = make(map[int]int)
	for _, rightNum := range rightNums {
		freqMap[rightNum]++
	}
	for _, leftNum := range leftNums {
		if freqMap[leftNum] > 0 {
			sscore += freqMap[leftNum] * leftNum
		}
	}
	pl("Similarity score:", sscore)
}
