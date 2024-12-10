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

type Tile struct {
	loc      Coords
	height   int
	trails   int
	visitors map[int]bool
}

type Map struct {
	width  int
	height int
	size   int
	grid   []Tile
}

func findTrailFrom(t Tile, m Map, visitor int) {
	if t.height == 9 {
		index := coordsToIndex(t.loc, m.width)
		m.grid[index].trails = m.grid[index].trails + 1
		m.grid[index].visitors[visitor] = true
		return
	}

	for _, neighbour := range []Coords{
		{t.loc.x - 1, t.loc.y},
		{t.loc.x + 1, t.loc.y},
		{t.loc.x, t.loc.y - 1},
		{t.loc.x, t.loc.y + 1}} {

		// ignore tiles outside the map
		if !isWithinBounds(neighbour, &m) {
			continue
		}

		neighbourTile := m.grid[coordsToIndex(neighbour, m.width)]
		if neighbourTile.height == t.height+1 {
			findTrailFrom(neighbourTile, m, visitor)
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

	// part 2 & 1
	sumP1 := 0
	sumP2 := 0

	// build the topological map
	mapSize := maxW * maxY
	topoMap := Map{maxW, maxY, mapSize, []Tile{}}
	indexToCoords := initIndexToCoords(topoMap.width)
	startingLocations := []int{}
	targetLocations := []int{}
	for i, c := range mapString {
		elevation := int(c - '0')
		tile := Tile{indexToCoords(i), elevation, 0, map[int]bool{}}
		topoMap.grid = append(topoMap.grid, tile)
		if elevation == 0 {
			startingLocations = append(startingLocations, i)
		} else if elevation == 9 {
			targetLocations = append(targetLocations, i)
		}
	}

	for _, currentStart := range startingLocations {
		findTrailFrom(topoMap.grid[currentStart], topoMap, currentStart)
	}

	for _, currentTarget := range targetLocations {
		for v, _ := range topoMap.grid[currentTarget].visitors {
			topoMap.grid[v].trails++
		}
		sumP2 += topoMap.grid[currentTarget].trails
	}

	for _, currentStart := range startingLocations {
		if topoMap.grid[currentStart].trails > 0 {
			sumP1 += topoMap.grid[currentStart].trails
		}
	}

	durationOverAll := time.Since(startOverAll)

	pl("Part 1:", sumP1)
	pl("Part 2:", sumP2)
	pl("Overall duration:", durationOverAll)
}
