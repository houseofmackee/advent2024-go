package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func initIndexToCoords(w int) func(index int) Coords {
	var cache = make(map[int]Coords)
	var width = w
	return func(index int) Coords {
		if c, ok := cache[index]; ok {
			return c
		}
		x := index % width
		y := index / width
		c := Coords{x, y}
		cache[index] = c
		return c
	}
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
	startOverAll := time.Now()

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
	cityMap := City{maxW, maxY, mapSize, []byte(mapString)}
	antiMap := City{maxW, maxY, mapSize, slices.Repeat([]byte{'.'}, mapSize)}
	indexToCoords := initIndexToCoords(cityMap.width)

	// Part 1
	startP1 := time.Now()
	for i, char := range cityMap.grid {
		// ignore empty spaces and or frequencies that only appear once
		if char == '.' { //} || count(cityMap.grid, char) == 1 {
			continue
		}

		// get the coordinates of the current antenna node
		freqCoords := indexToCoords(i)

		// find all other nodes with the same frequency
		for j := i + 1; j < cityMap.size; j++ {
			if cityMap.grid[j] == char {
				// get the coordinates of the paired antenna node
				pairCoords := indexToCoords(j)

				// calculate the distance between the two nodes
				diff := coordsDiff(freqCoords, pairCoords)

				// put the anti nodes in the antiMap on opposite sides of the two nodes
				antiA := Coords{freqCoords.x + diff.x, freqCoords.y + diff.y}
				antiB := Coords{pairCoords.x - diff.x, pairCoords.y - diff.y}

				putAntiNode(antiA, &antiMap)
				putAntiNode(antiB, &antiMap)
			}
		}
	}

	sumP1 := count(antiMap.grid, '#')
	durationP1 := time.Since(startP1)

	// Part 2
	startP2 := time.Now()
	for i, char := range cityMap.grid {
		// ignore empty spaces and or frequencies that only appear once
		if char == '.' { //|| count(cityMap.grid, char) == 1 {
			continue
		}

		// get the coordinates of the current antenna node
		freqCoords := indexToCoords(i)

		// find all other nodes with the same frequency
		for j := i + 1; j < cityMap.size; j++ {
			if cityMap.grid[j] == char {
				// get the coordinates of the paired antenna node
				pairCoords := indexToCoords(j)

				// calculate the distance between the two nodes
				diff := coordsDiff(freqCoords, pairCoords)

				// draw a dotted line extending both ways from the two nodes
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

	sumP2 := count(antiMap.grid, '#')
	durationP2 := time.Since(startP2)
	durationOverAll := time.Since(startOverAll)

	pl("Part 1:", sumP1)
	pl("Part 1 duration:", durationP1)

	pl("Part 2:", sumP2)
	pl("Part 2 duration:", durationP2)

	pl("Overall duration:", durationOverAll)

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
