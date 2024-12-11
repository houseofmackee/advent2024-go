package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"
	"time"
)

// I am lazy...
var pl = fmt.Println

// function to cut a string in half
func splitStringInHalf(s string) (left, right string) {
	return s[:len(s)/2], s[len(s)/2:]
}

// split a number in half
func initSplitNumberInHalf() func(n int) (left, right int) {
	var cache = make(map[int]struct {
		left  int
		right int
	})
	return func(n int) (int, int) {
		if v, ok := cache[n]; ok {
			return v.left, v.right
		}
		a, b := splitStringInHalf(strconv.Itoa(n))
		l, _ := strconv.Atoi(a)
		r, _ := strconv.Atoi(b)
		cache[n] = struct {
			left  int
			right int
		}{l, r}
		return l, r
	}
}

func initIsEvenLength() func(n int) bool {
	var cache = make(map[int]bool)
	return func(n int) bool {
		if v, ok := cache[n]; ok {
			return v
		}
		s := strconv.Itoa(n)
		cache[n] = len(s)%2 == 0
		return cache[n]
	}
}

type Number interface {
	int | float32 | float64
}

func sumValues[K comparable, V Number](m *map[K]V) V {
	var sum V
	for _, v := range *m {
		sum += v
	}
	return sum
}

var splitNumberInHalf = initSplitNumberInHalf()
var isEvenLength = initIsEvenLength()

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
	var inputString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		inputString += strings.TrimSpace(line)
	}

	// part 1 & 2
	p1Start := time.Now()
	p2Start := time.Now()
	var p1End time.Time
	var p2End time.Time

	sumP1 := 0
	sumP2 := 0
	numP1Blinks := 25
	numP2Blinks := 75

	// turn string of numbers into slice of values and count the occurences in a map
	values := strings.Split(inputString, " ")
	numMap := make(map[int]int)
	for _, value := range values {
		v, _ := strconv.Atoi(value)
		numMap[v]++
	}

	blinks := 0
	for {
		newNums := make(map[int]int)
		for k, v := range numMap {
			delete(numMap, k)
			if k == 0 {
				newNums[1] += v
			} else if isEvenLength(k) {
				a, b := splitNumberInHalf(k)
				newNums[a] += v
				newNums[b] += v
			} else {
				newNums[k*2024] += v
			}
		}
		numMap = maps.Clone(newNums)

		// part 1
		blinks++
		if blinks == numP1Blinks {
			sumP1 = sumValues(&numMap)
			p1End = time.Now()
		}

		// part 2
		if blinks == numP2Blinks {
			sumP2 = sumValues(&numMap)
			p2End = time.Now()
			break
		}
	}

	durationOverAll := time.Since(startOverAll)
	p1Duration := p1End.Sub(p1Start)
	p2Duration := p2End.Sub(p2Start)

	pl("Part 1:", sumP1)
	pl("Duration Part 1:", p1Duration)
	pl("Part 2:", sumP2)
	pl("Duration Part 2:", p2Duration)
	pl("Overall duration:", durationOverAll)
}
