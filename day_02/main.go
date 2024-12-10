package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// function to remove an element from a slice
func removeInt(slice []int, i int) []int {
	var newSlice []int
	for j := 0; j < len(slice); j++ {
		if j != i {
			newSlice = append(newSlice, slice[j])
		}
	}
	return newSlice
}

// function to check if a report is safe
func isReportSafe(intValues []int) bool {
	up, down := false, false
	for i := 0; i < len(intValues)-1; i++ {
		left, right := intValues[i], intValues[i+1]
		diff := absDiffInt(left, right)
		if diff == 0 || diff > 3 {
			return false
		}

		if left < right {
			up = true
		} else {
			down = true
		}

		if up && down {
			return false
		}
	}
	return true
}

func main() {

	// open the file called input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		pl("Error: ", err)
		return
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var p1SafeReports, p1UnsafeReports int
	var p2SafeReports, p2UnsafeReports int
	for fileScanner.Scan() {
		var intValues []int
		values := strings.Split(fileScanner.Text(), " ")
		for _, v := range values {
			intValue, _ := strconv.Atoi(v)
			intValues = append(intValues, intValue)
		}
		pl("Int values:", intValues)

		// check for Part 1
		if isReportSafe(intValues) {
			p1SafeReports++
		} else {
			p1UnsafeReports++
		}

		// check for Part 2
		if len(intValues) > 2 {
			isSafe := false
			for i := 0; i < len(intValues); i++ {
				tempIntValues := removeInt(intValues, i)
				if isReportSafe(tempIntValues) {
					isSafe = true
					break
				}
			}
			if isSafe {
				p2SafeReports++
			} else {
				p2UnsafeReports++
			}
		} else {
			pl("Not enough values to check for Part 2")
		}
	}

	pl("Part 1")
	pl("Safe reports:", p1SafeReports)
	pl("Unsafe reports:", p1UnsafeReports, "\n")

	pl("Part 2")
	pl("Safe reports:", p2SafeReports)
	pl("Unsafe reports:", p2UnsafeReports)
}
