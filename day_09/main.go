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

type BlockInfo struct {
	id   int
	size int
}

func findLast[T comparable](slice []T, target T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == target {
			return i
		}
	}
	return -1
}

func findLastNot[T comparable](slice []T, target T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] != target {
			return i
		}
	}
	return -1
}

func findFirstFree(slice []int) int {
	for i, v := range slice {
		if v == -1 {
			return i
		}
	}
	return -1
}

func findFirstFreeN(slice []int, n int) int {
	for i := 0; i < len(slice)-n; i++ {
		if slice[i] == -1 {
			found := true
			for j := 1; j < n; j++ {
				if slice[i+j] != -1 {
					found = false
					break
				}
			}
			if found {
				return i
			}
		}
	}
	return -1
}

func checksum(slice []int) int {
	checksum := 0
	for i, v := range slice {
		if v != -1 {
			checksum += v * i
		}
	}
	return checksum
}

func swap(slice []int, a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
}

func swapN(slice []int, a, b, n int) {
	for i := 0; i < n; i++ {
		swap(slice, a+i, b+i)
	}
}

func count[T comparable](values []T, value T) int {
	count := 0
	for _, item := range values {
		if item == value {
			count++
		}
	}
	return count
}

func countFrom[T comparable](values []T, value T, start int) int {
	count := 0
	for i := start; i < len(values); i++ {
		if values[i] == value {
			count++
		}
	}
	return count
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
	var diskString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		diskString += line
	}
	diskString = strings.TrimSpace(diskString)

	////////////////////////////////////////
	// part 1
	startP1 := time.Now()

	// turn digits string into slice of bytes with the same values
	numBlocks := len(diskString)
	diskBytes := make([]byte, numBlocks)
	for i, c := range diskString {
		diskBytes[i] = byte(c - '0')
	}

	// build block map with sizes and IDs
	diskMap := make([]BlockInfo, numBlocks)
	currentId := 0
	for i := 0; i < numBlocks; i++ {
		// check if disk block is a file
		if i%2 == 0 {
			diskMap[i] = BlockInfo{currentId, int(diskBytes[i])}
			currentId++
		} else {
			// process empty blocks
			diskMap[i] = BlockInfo{-1, int(diskBytes[i])}
		}
	}

	// build slice with IDs of the blocks
	blockMap := []int{}
	for i := 0; i < numBlocks; i++ {
		if diskMap[i].id == -1 {
			for j := 0; j < diskMap[i].size; j++ {
				blockMap = append(blockMap, -1)
			}
		} else {
			id := diskMap[i].id
			for j := 0; j < diskMap[i].size; j++ {
				blockMap = append(blockMap, id)
			}
		}
	}

	// for use in part 2
	defragMap := slices.Clone(blockMap)

	for {
		firstFree := slices.Index(blockMap, -1)
		if firstFree == -1 {
			break
		}

		lastUsed := findLastNot(blockMap, -1)
		if lastUsed == -1 || firstFree > lastUsed {
			break
		}

		// swap the blocks
		swap(blockMap, firstFree, lastUsed)
	}

	// calculate checksum
	sumP1 := checksum(blockMap)
	durationP1 := time.Since(startP1)

	////////////////////////////////////////
	// part 2
	startP2 := time.Now()
	for nextId := currentId - 1; nextId >= 0; nextId-- {
		fileOffset := slices.Index(defragMap, nextId)
		fileLength := countFrom(defragMap, nextId, fileOffset)
		firstFree := findFirstFreeN(defragMap, fileLength)

		// ignore if there's no free space or if free space is after the file
		if firstFree == -1 || firstFree > fileOffset {
			continue
		}

		// move the file to the free space
		swapN(defragMap, firstFree, fileOffset, fileLength)
	}

	// calculate checksum
	sumP2 := checksum(defragMap)
	durationP2 := time.Since(startP2)
	durationOverAll := time.Since(startOverAll)

	pl("Part 1:", sumP1)
	pl("Part 1 duration:", durationP1)

	pl("Part 2:", sumP2)
	pl("Part 2 duration:", durationP2)

	pl("Overall duration:", durationOverAll)
}
