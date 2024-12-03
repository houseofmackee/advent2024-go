package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// I am lazy...
var pl = fmt.Println

// function to process all 'mul' patterns in a string and return the sum
func muls(input string) (int, error) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("no valid 'mul' patterns found in input string")
	}

	var result int
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		num1, err1 := strconv.Atoi(match[1])
		num2, err2 := strconv.Atoi(match[2])

		if err1 != nil || err2 != nil {
			continue
		}

		result += num1 * num2
	}
	return result, nil
}

// function to process a line of input and return the sum of all 'mul' patterns
// and obey the 'do' and 'don't' commands
func processLine(input string, blocker *bool) (int, error) {
	re := regexp.MustCompile(`do\(\)|don't\(\)|mul\(\d{1,3},\d{1,3}\)`)
	matches := re.FindAllString(input, -1)

	var result int
	for _, match := range matches {
		if match == "do()" {
			*blocker = false
		} else if match == "don't()" {
			*blocker = true
		} else if !*blocker {
			value, _ := muls(match)
			result += value
		}
	}
	return result, nil
}

func main() {

	// open the file called input.txt
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
	}
	defer file.Close()

	// create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// read the file
	var p1Result int
	var p2Result int
	var blocker = false
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			pl(err)
		}

		// process Part 1
		value, err := muls(line)
		if err != nil {
			pl(err)
		}
		p1Result += value

		// process Part 2
		value, err = processLine(line, &blocker)
		if err != nil {
			pl(err)
		}
		p2Result += value
	}
	pl("Part 1: ", p1Result)
	pl("Part 2: ", p2Result)
}
