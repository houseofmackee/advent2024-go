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

type Guard struct {
	x      int
	y      int
	facing int
}

type Room struct {
	width  int
	height int
	grid   string
}

func indexToCoords(index, width int) (int, int) {
	x := index % width
	y := index / width
	return x, y
}

func coordsToIndex(x, y, width int) int {
	return y*width + x
}

func checkRoomLocation(x, y int, room *Room) string {
	if x < 0 || x >= room.width || y < 0 || y >= room.height {
		return "offgrid"
	}

	if room.grid[coordsToIndex(x, y, room.width)] == '#' {
		return "blocked"
	}

	return "valid"
}

func moveGuardInRoom(guard *Guard, room *Room) string {
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
	if move == "offgrid" {
		return move
	} else if move == "valid" {
		guard.x += x
		guard.y += y
		room.grid = room.grid[:coordsToIndex(guard.x, guard.y, room.width)] + "*" + room.grid[coordsToIndex(guard.x, guard.y, room.width)+1:]
	} else if move == "blocked" {
		guard.facing++
		if guard.facing > 3 {
			guard.facing = 0
		}
	}
	return move
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
	}
	defer file.Close()

	var maxW, maxY int
	var grid string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maxW = len(line)
		maxY++
		grid += line
	}
	grid = strings.Replace(grid, "^", "X", -1)
	room := Room{maxW, maxY, grid}

	guardIndex := strings.Index(grid, "X")
	posX, posY := indexToCoords(guardIndex, maxW)
	guard := Guard{posX, posY, 0}

	// process Part 1
	for {
		if moveGuardInRoom(&guard, &room) == "offgrid" {
			break
		}
	}
	pl("Part 1:", strings.Count(room.grid, "*"))

	// process Part 2
	grid = room.grid
	loops := 0
	start := time.Now()
	for i := 0; i < len(grid); i++ {
		if grid[i] != '*' || i == guardIndex {
			continue
		}

		tempGrid := grid[:i] + "#" + grid[i+1:]
		room = Room{maxW, maxY, tempGrid}
		guard := Guard{posX, posY, 0}
		guardHistory := []Guard{}

		for {
			moveResult := moveGuardInRoom(&guard, &room)
			if moveResult == "offgrid" {
				break
			}
			if moveResult == "valid" {
				continue
			}
			if moveResult == "blocked" {
				for j := 0; j < len(guardHistory); j++ {
					if guardHistory[j] == guard {
						moveResult = "loop"
						break
					}
				}
			}
			if moveResult != "loop" {
				guardHistory = append(guardHistory, guard)
			} else {
				loops++
				break
			}
		}
	}
	pl("Time:", time.Since(start))
	pl("Part 2:", loops)
}
