package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// I am lazy...
var pl = fmt.Println

type Coords struct {
	x int
	y int
}

type City struct {
	width  int
	height int
	size   int
	grid   []byte
}

// function to get the absolute difference betwen two ints
func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func indexToCoords(index, width int) Coords {
	x := index % width
	y := index / width
	return Coords{x, y}
}

func isWithinBounds(c Coords, city *City) bool {
	return c.x >= 0 && c.x < city.width && c.y >= 0 && c.y < city.height
}

func putAntiNode(c Coords, city *City) bool {
	if isWithinBounds(c, city) {
		city.grid[coordsToIndex(c, city.width)] = '#'
		return true
	}
	return false
}

func coordsToIndex(c Coords, width int) int {
	return c.y*width + c.x
}

func coordsDiff(c1, c2 Coords) Coords {
	return Coords{c1.x - c2.x, c1.y - c2.y}
}

func coordsDiffAbs(c1, c2 Coords) Coords {
	return Coords{absDiffInt(c1.x, c2.x), absDiffInt(c1.y, c2.y)}
}

func count(grid []byte, value byte) int {
	count := 0
	for _, char := range grid {
		if char == value {
			count++
		}
	}
	return count
}

func main() {
	// Open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
		return
	}
	defer file.Close()

	// Parse the input file
	var maxW, maxY int
	var mapString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maxW = len(line)
		maxY++
		mapString += line
	}

	mapSize := maxW * maxY
	mapGrid := []byte(mapString)
	antiGrid := make([]byte, mapSize)

	cityMap := City{maxW, maxY, mapSize, mapGrid}
	antiMap := City{maxW, maxY, mapSize, antiGrid}
	// fill the antiMap with empty spaces
	for i := 0; i < antiMap.size; i++ {
		antiMap.grid[i] = '.'
	}

	// Part 1
	start := time.Now()
	for i, char := range cityMap.grid {
		// ignore empty spaces and or frequencies that only appear once
		if char == '.' { //} || count(cityMap.grid, char) == 1 {
			continue
		}
		freqCoords := indexToCoords(i, cityMap.width)

		// find all other nodes with the same frequency
		for j := i + 1; j < cityMap.size; j++ {
			if cityMap.grid[j] == char {
				pairCoords := indexToCoords(j, cityMap.width)

				// calculate the distance between the two nodes
				diff := coordsDiff(freqCoords, pairCoords)
				antiA := Coords{freqCoords.x + diff.x, freqCoords.y + diff.y}
				antiB := Coords{pairCoords.x - diff.x, pairCoords.y - diff.y}

				putAntiNode(antiA, &antiMap)
				putAntiNode(antiB, &antiMap)
			}
		}
	}

	sum := count(antiMap.grid, '#')
	durartion := time.Since(start)
	pl("Part 1:", sum)
	pl("Part 1 duration:", durartion)

	// Part 2
	start = time.Now()
	for i, char := range cityMap.grid {
		// ignore empty spaces and or frequencies that only appear once
		if char == '.' { //|| count(cityMap.grid, char) == 1 {
			continue
		}
		freqCoords := indexToCoords(i, cityMap.width)

		// find all other nodes with the same frequency
		for j := i + 1; j < cityMap.size; j++ {
			if cityMap.grid[j] == char {
				pairCoords := indexToCoords(j, cityMap.width)

				// calculate the distance between the two nodes
				diff := coordsDiff(freqCoords, pairCoords)
				lineLen := 1
				for {
					antiA := Coords{freqCoords.x + (lineLen * diff.x), freqCoords.y + (lineLen * diff.y)}
					if !putAntiNode(antiA, &antiMap) {
						break
					}
					lineLen++
				}

				lineLen = 1
				for {
					antiB := Coords{pairCoords.x - (lineLen * diff.x), pairCoords.y - (lineLen * diff.y)}
					if !putAntiNode(antiB, &antiMap) {
						break
					}
					lineLen++
				}
			}
		}
		// turn antenna node into an anti node
		putAntiNode(freqCoords, &antiMap)
	}

	sum = count(antiMap.grid, '#')
	durartion = time.Since(start)
	pl("Part 2:", sum)
	pl("Part 2 duration:", durartion)

	// for i := 0; i < antiMap.size; i++ {
	// 	if i%antiMap.width == 0 {
	// 		pl()
	// 	}
	// 	if antiMap.grid[i] == '#' {
	// 		fmt.Printf("%c", antiMap.grid[i])
	// 	} else {
	// 		fmt.Printf("%c", cityMap.grid[i])
	// 	}
	// }
	// pl()
}
