package main

import (
	"bufio"
	"fmt"
	"os"
)

// I am lazy...
var pl = fmt.Println

func matchXmas(x, m, a, s byte) bool {
	result :=
		(string(x) == "X" && string(m) == "M" && string(a) == "A" && string(s) == "S") ||
			(string(x) == "S" && string(m) == "A" && string(a) == "M" && string(s) == "X")
	return result
}

func matchMas(m, a, s byte) bool {
	result :=
		(string(m) == "M" && string(a) == "A" && string(s) == "S") ||
			(string(m) == "S" && string(a) == "A" && string(s) == "M")
	return result
}

func findXmas(lines []string, x, y int) int {

	count := 0
	maxX := len(lines[0]) - 1
	maxY := len(lines) - 1

	//   XMAS
	//  MMM
	// A A A
	//S  S  S

	if (x < maxX-2) && matchXmas(
		lines[y][x],
		lines[y][x+1],
		lines[y][x+2],
		lines[y][x+3]) {
		count++
	}

	if (x < maxX-2) && (y < maxY-2) && matchXmas(
		lines[y][x],
		lines[y+1][x+1],
		lines[y+2][x+2],
		lines[y+3][x+3]) {
		count++
	}

	if (y < maxY-2) && matchXmas(
		lines[y][x],
		lines[y+1][x],
		lines[y+2][x],
		lines[y+3][x]) {
		count++
	}

	if (x > 2) && (y < maxY-2) && matchXmas(
		lines[y][x],
		lines[y+1][x-1],
		lines[y+2][x-2],
		lines[y+3][x-3]) {
		count++
	}

	return count
}

func findMas(lines []string, x, y int) int {
	count := 0
	if matchMas(
		lines[y+0][x+0],
		lines[y+1][x+1],
		lines[y+2][x+2]) &&
		matchMas(
			lines[y+0][x+2],
			lines[y+1][x+1],
			lines[y+2][x+0]) {
		count++
	}
	return count
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error reading file")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Process part 1
	count := 0
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			count += findXmas(lines, x, y)
		}
	}
	pl("Part 1:", count)

	// Process part 2
	count = 0
	for y := 0; y < len(lines)-2; y++ {
		for x := 0; x < len(lines[y])-2; x++ {
			count += findMas(lines, x, y)
		}
	}
	pl("Part 2:", count)
}
