package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// I am lazy...
var pl = fmt.Println

func checkRule(rule, update string) string {
	rules := strings.Split(rule, "|")
	indexA := strings.Index(update, rules[0])
	indexB := strings.Index(update, rules[1])

	if indexA < 0 || indexB < 0 {
		return "ignore"
	}

	if indexA > indexB {
		return "invalid"
	}

	return "valid"
}

func applyRule(rule, update string) string {
	rules := strings.Split(rule, "|")
	indexA := strings.Index(update, rules[0])
	indexB := strings.Index(update, rules[1])

	if indexA < 0 || indexB < 0 {
		return update
	}

	if indexA < indexB {
		return update
	}

	update = strings.Replace(update, rules[0], "replaceme", 1)
	update = strings.Replace(update, rules[1], rules[0], 1)
	update = strings.Replace(update, "replaceme", rules[1], 1)

	return update
}

func getMiddleNum(update string) int {
	pages := strings.Split(update, ",")
	middle := len(pages) / 2
	middleNum, _ := strconv.Atoi(pages[middle])
	return middleNum
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		pl(err)
	}
	defer file.Close()

	var rules []string
	var updates []string
	var readingRules = true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingRules = false
			continue
		}
		if readingRules {
			rules = append(rules, line)
		} else {
			updates = append(updates, line)
		}
	}

	// process Part 1
	middleSum := 0
	var invalids []string
	for i := 0; i < len(updates); i++ {
		update := updates[i]
		isValid := true
		for j := 0; j < len(rules); j++ {
			ruleStatus := checkRule(rules[j], update)
			if ruleStatus == "ignore" {
				continue
			}
			if ruleStatus == "invalid" {
				invalids = append(invalids, update)
				isValid = false
				break
			}
		}
		if isValid {
			middleSum += getMiddleNum(update)
		}
	}
	pl("Part 1:", middleSum)

	// process Part 2
	middleSum = 0
	for i := 0; i < len(invalids); i++ {
		invalid := invalids[i]

		for {
			needRerun := false
			for j := 0; j < len(rules); j++ {
				new := applyRule(rules[j], invalid)
				if new != invalid {
					invalid = new
					needRerun = true
				}
			}
			if !needRerun {
				break
			}
		}
		middleSum += getMiddleNum(invalid)
	}
	pl("Part 2:", middleSum)
}
