package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// I am lazy...
var pl = fmt.Println

type Coords struct {
	x int
	y int
}

type Borders struct {
	count  int
	facing map[int]bool
}

type Tile struct {
	loc     Coords
	group   int
	borders Borders
	id      int
}

type Map struct {
	width  int
	height int
	size   int
	grid   []Tile
}

const (
	// directions
	NORTH = 2
	SOUTH = 3
	WEST  = 0
	EAST  = 1
)

func findRegion(m *Map, plotIndex, plotId int) {
	if m.grid[plotIndex].id >= 0 {
		return
	}
	m.grid[plotIndex].id = plotId

	for i, neighbour := range []Coords{
		{m.grid[plotIndex].loc.x - 1, m.grid[plotIndex].loc.y},
		{m.grid[plotIndex].loc.x + 1, m.grid[plotIndex].loc.y},
		{m.grid[plotIndex].loc.x, m.grid[plotIndex].loc.y - 1},
		{m.grid[plotIndex].loc.x, m.grid[plotIndex].loc.y + 1}} {

		// ignore tiles outside the map
		if !isWithinBounds(neighbour, m) {
			m.grid[plotIndex].borders.count++
			m.grid[plotIndex].borders.facing[i] = true
			continue
		}
		j := coordsToIndex(neighbour, m.width)
		if m.grid[j].group == m.grid[plotIndex].group {
			findRegion(m, j, plotId)
		} else {
			m.grid[plotIndex].borders.count++
			m.grid[plotIndex].borders.facing[i] = true
		}
	}
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

func isWithinBounds(c Coords, m *Map) bool {
	return c.x >= 0 && c.x < m.width && c.y >= 0 && c.y < m.height
}

func coordsToIndex(c Coords, width int) int {
	return c.y*width + c.x
}

func main() {
	startOverAll := time.Now()

	// open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
		return
	}
	defer file.Close()

	// parse the input file
	var maxW, maxY int
	var mapString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		maxW = len(line)
		maxY++
		mapString += strings.TrimSpace(line)
	}

	// part 1
	sumP1 := 0
	sumP2 := 0
	startP1 := time.Now()

	// build the topological map
	mapSize := maxW * maxY
	topoMap := Map{maxW, maxY, mapSize, []Tile{}}
	indexToCoords := initIndexToCoords(topoMap.width)
	for i, c := range mapString {
		plantType := int(c - 'A')
		tile := Tile{
			loc:     indexToCoords(i),
			group:   plantType,
			borders: Borders{count: 0, facing: make(map[int]bool)},
			id:      -1}
		topoMap.grid = append(topoMap.grid, tile)
	}

	// find regions
	plotId := 0
	regionsGroup := make(map[int]int)
	for i := 0; i < mapSize; i++ {
		if topoMap.grid[i].id < 0 {
			findRegion(&topoMap, i, plotId)
			regionsGroup[plotId] = topoMap.grid[i].group
			plotId++
		}
	}

	// find number of unique regions
	uniqueRegions := make(map[int]bool)
	for i := 0; i < mapSize; i++ {
		uniqueRegions[topoMap.grid[i].id] = true
	}

	borders := make(map[int]int)
	for i := 0; i < mapSize; i++ {
		if topoMap.grid[i].borders.count > 0 {
			borders[topoMap.grid[i].id] += topoMap.grid[i].borders.count
		}
	}

	areas := make(map[int]int)
	for i := 0; i < mapSize; i++ {
		areas[topoMap.grid[i].id]++
	}

	// print regions with borders and areas
	for k := range uniqueRegions {
		sumP1 += areas[k] * borders[k]
	}
	endP1 := time.Since(startP1)

	// part 2
	startP2 := time.Now()

	// find all straight, uniterrupted, lines facing N/E/S/W
	regionSides := make(map[int]int)
	activeBorder := false

	// helper function to walk a line and tell if it started or ended
	walkLine := func(x, y, r, direction int) bool {
		i := coordsToIndex(Coords{x, y}, topoMap.width)
		if topoMap.grid[i].id == r {
			if topoMap.grid[i].borders.facing[direction] {
				if !activeBorder {
					regionSides[r]++
					activeBorder = true
				}
			} else {
				activeBorder = false
			}
		} else {
			activeBorder = false
		}
		return activeBorder
	}

	// check all regions
	for r := range uniqueRegions {

		// find MORTH and SOUTH facing borders
		for y := 0; y < topoMap.height; y++ {

			// find NORTH facing borders
			activeBorder = false
			for x := 0; x < topoMap.width; x++ {
				activeBorder = walkLine(x, y, r, NORTH)
			}

			// find SOUTH facing borders
			activeBorder = false
			for x := 0; x < topoMap.width; x++ {
				activeBorder = walkLine(x, y, r, SOUTH)
			}
		}

		// find EAST and WEST facing borders
		for x := 0; x < topoMap.width; x++ {

			// find EAST facing borders
			activeBorder = false
			for y := 0; y < topoMap.height; y++ {
				activeBorder = walkLine(x, y, r, EAST)
			}

			// find WEST facing borders
			activeBorder = false
			for y := 0; y < topoMap.height; y++ {
				activeBorder = walkLine(x, y, r, WEST)
			}
		}
	}

	for k := range uniqueRegions {
		sumP2 += regionSides[k] * areas[k]
	}

	endP2 := time.Since(startP2)
	durationOverAll := time.Since(startOverAll)

	pl("Part 1:", sumP1)
	pl("Duration part 1:", endP1)
	pl("Part 2:", sumP2)
	pl("Duration part 2:", endP2)
	pl("Overall duration:", durationOverAll)
}
