package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

// I am lazy...
var pl = fmt.Println

const (
	OFFGRID = iota
	BLOCKED
	VALID
	LOOP
)

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

type Guard struct {
	x      int
	y      int
	facing int
}

type Room struct {
	width  int
	height int
	grid   []byte
}

func indexToCoords(index, width int) (int, int) {
	x := index % width
	y := index / width
	return x, y
}

func coordsToIndex(x, y, width int) int {
	return y*width + x
}

func checkRoomLocation(x, y int, room *Room) int {
	if x < 0 || x >= room.width || y < 0 || y >= room.height {
		return OFFGRID
	} else if room.grid[coordsToIndex(x, y, room.width)] == '#' {
		return BLOCKED
	}
	return VALID
}

func moveGuardInRoom(guard *Guard, room *Room) int {
	var x, y int
	switch guard.facing {
	case NORTH:
		x, y = 0, -1
	case EAST:
		x, y = 1, 0
	case SOUTH:
		x, y = 0, 1
	case WEST:
		x, y = -1, 0
	}

	x += guard.x
	y += guard.y
	move := checkRoomLocation(x, y, room)
	if move == VALID {
		guard.x = x
		guard.y = y
		room.grid[coordsToIndex(guard.x, guard.y, room.width)] = '*'
	} else if move == BLOCKED {
		guard.facing++
		if guard.facing > 3 {
			guard.facing = 0
		}
	}
	return move
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
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
	}
	defer file.Close()

	var maxW, maxY int
	var grids string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maxW = len(line)
		maxY++
		grids += line
	}
	grids = strings.Replace(grids, "^", "X", -1)
	guardIndex := strings.Index(grids, "X")
	grid := []byte(grids)
	room := Room{maxW, maxY, grid}

	posX, posY := indexToCoords(guardIndex, maxW)
	guard := Guard{posX, posY, NORTH}

	// process Part 1
	for {
		if moveGuardInRoom(&guard, &room) == OFFGRID {
			break
		}
	}
	pl("Part 1:", count(room.grid, '*'))

	// process Part 2
	mainGrid := slices.Clone(room.grid)
	mainGridLen := len(mainGrid)
	loops := 0
	start := time.Now()
	for i := 0; i < mainGridLen; i++ {
		if grid[i] != '*' || i == guardIndex {
			continue
		}

		tempGrid := slices.Clone(mainGrid)
		tempGrid[i] = '#'
		tempRoom := Room{maxW, maxY, tempGrid}
		tempGuard := Guard{posX, posY, NORTH}
		guardHistory := []Guard{}

		for {
			moveResult := moveGuardInRoom(&tempGuard, &tempRoom)
			if moveResult == VALID {
				continue
			}
			if moveResult == OFFGRID {
				break
			}
			if slices.Contains(guardHistory, tempGuard) {
				loops++
				break
			}
			guardHistory = append(guardHistory, tempGuard)
		}
	}
	pl("Time:", time.Since(start))
	pl("Part 2:", loops)
}
