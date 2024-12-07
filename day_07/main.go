package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// I am lazy...
var pl = fmt.Println

type Equation struct {
	result int
	values []int
}

func parseEquationString(input string) (Equation, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) < 2 {
		return Equation{}, fmt.Errorf("input does not contain a colon")
	}

	beforeStr := strings.TrimSpace(parts[0])
	beforeInt, err := strconv.Atoi(beforeStr)
	if err != nil {
		return Equation{}, fmt.Errorf("error parsing '%s' as integer: %v", beforeStr, err)
	}

	afterStr := strings.TrimSpace(parts[1])
	afterFields := strings.Fields(afterStr)

	var afterInts []int
	for _, field := range afterFields {
		val, err := strconv.Atoi(field)
		if err != nil {
			return Equation{}, fmt.Errorf("error parsing '%s' as integer: %v", field, err)
		}
		afterInts = append(afterInts, val)
	}

	return Equation{beforeInt, afterInts}, nil
}

// function type for generating all possible combinations of operators
type RuneGenerator func(int) [][]rune

// generate all possible combinations of + and * operators
func initGenerateCombinations() func(length int) [][]rune {
	var cache = make(map[int][][]rune)

	return func(length int) [][]rune {
		if combo, ok := cache[length]; ok {
			return combo
		}
		// Total number of combinations is 2^length
		total := 1 << length
		result := make([][]rune, 0, total)
		for i := 0; i < total; i++ {
			combo := make([]rune, length)
			for pos := 0; pos < length; pos++ {
				if (i & (1 << pos)) == 0 {
					combo[pos] = '*'
				} else {
					combo[pos] = '+'
				}
			}
			result = append(result, combo)
		}
		cache[length] = result
		return result
	}
}

// generate all possible combinations of +, *, and | operators
func initGenerateConcatCombinations() func(length int) [][]rune {
	var cache = make(map[int][][]rune)

	return func(length int) [][]rune {
		if combo, ok := cache[length]; ok {
			return combo
		}
		chars := []rune{'*', '|', '+'}
		total := int(math.Pow(3, float64(length))) // 3^length combinations
		result := make([][]rune, 0, total)
		for i := 0; i < total; i++ {
			combo := make([]rune, length)
			value := i
			for pos := 0; pos < length; pos++ {
				digit := value % 3
				value = value / 3
				combo[pos] = chars[digit]
			}
			result = append(result, combo)
		}
		cache[length] = result
		return result
	}
}

func concatenate(a, b int) int {
	strA := strconv.Itoa(a)
	strB := strconv.Itoa(b)
	combined, _ := strconv.Atoi(strA + strB)
	return combined
}

// brute force all possible combinations of operators
func processEquation(equation Equation, gen RuneGenerator) int {
	target := equation.result
	values := equation.values
	numOps := len(values)
	runes := gen(numOps - 1)

	for _, combo := range runes {
		running := values[0]
		for i := 0; i < numOps-1; i++ {
			if combo[i] == '|' {
				running = concatenate(running, values[i+1])
			} else if combo[i] == '*' {
				running *= values[i+1]
			} else {
				running += values[i+1]
			}
			if running > target {
				break
			}
		}
		if running == target {
			return running
		}
	}

	return 0
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
	}
	defer file.Close()

	equations := []Equation{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		equation, err := parseEquationString(line)
		if err != nil {
			pl(err)
			continue
		}
		equations = append(equations, equation)
	}

	// process Part 1
	start := time.Now()
	sump1 := 0
	invalids := []Equation{}
	generateCombinations := initGenerateCombinations()
	for _, equation := range equations {
		result := processEquation(equation, generateCombinations)
		if result == 0 {
			invalids = append(invalids, equation)
		} else {
			sump1 += result
		}
	}
	pl("Time:", time.Since(start))
	pl("Part 1:", sump1)

	// process Part 2
	start = time.Now()
	sump2 := 0
	generateCombinations = initGenerateConcatCombinations()
	for _, equation := range invalids {
		sump2 += processEquation(equation, generateCombinations)
	}
	pl("Time:", time.Since(start))
	pl("Part 2:", sump1+sump2)
}
