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

const OFFGRID = 0
const BLOCKED = 1
const VALID = 2
const LOOP = 3

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
	}

	if room.grid[coordsToIndex(x, y, room.width)] == '#' {
		return BLOCKED
	}

	return VALID
}

func moveGuardInRoom(guard *Guard, room *Room) int {
	var x, y int
	switch guard.facing {
	case 0:
		x, y = 0, -1
	case 1:
		x, y = 1, 0
	case 2:
		x, y = 0, 1
	case 3:
		x, y = -1, 0
	}

	move := checkRoomLocation(guard.x+x, guard.y+y, room)
	if move == OFFGRID {
		return move
	} else if move == VALID {
		guard.x += x
		guard.y += y
		room.grid[coordsToIndex(guard.x, guard.y, room.width)] = '*' //room.grid[:coordsToIndex(guard.x, guard.y, room.width)] + "*" + room.grid[coordsToIndex(guard.x, guard.y, room.width)+1:]
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

func printRoom(room *Room) {
	for i, char := range room.grid {
		if i%room.width == 0 {
			pl()
		}
		fmt.Printf("%c", char)
	}
	pl()
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
	guard := Guard{posX, posY, 0}

	// process Part 1
	for {
		if moveGuardInRoom(&guard, &room) == OFFGRID {
			break
		}
	}
	pl("Part 1:", count(room.grid, '*'))

	// process Part 2
	mainGrid := room.grid
	mainGridLen := len(mainGrid)
	loops := 0
	start := time.Now()
	for i := 0; i < len(grid); i++ {
		if grid[i] != '*' || i == guardIndex {
			continue
		}

		tempGrid := make([]byte, mainGridLen)
		copy(tempGrid, mainGrid)
		tempGrid[i] = '#'
		tempRoom := Room{maxW, maxY, tempGrid}
		tempGuard := Guard{posX, posY, 0}
		guardHistory := []Guard{}

		for {
			moveResult := moveGuardInRoom(&tempGuard, &tempRoom)
			if moveResult == VALID {
				continue
			}
			if moveResult == OFFGRID {
				break
			}
			if moveResult == BLOCKED {

				for _, entry := range guardHistory {
					if entry == tempGuard {
						moveResult = LOOP
						break
					}
				}
			}
			if moveResult != LOOP {
				guardHistory = append(guardHistory, tempGuard)
			} else {
				loops++
				break
			}
		}
	}
	pl("Time:", time.Since(start))
	pl("Part 2:", loops)
}
