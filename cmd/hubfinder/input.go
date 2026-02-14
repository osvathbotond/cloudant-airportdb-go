package main

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func readFloatUntilValid(scanner *bufio.Scanner, variableName string, minValue, maxValue float64) float64 {
	for {
		fmt.Printf("Please enter the %s: ", variableName)
		if !scanner.Scan() {
			fmt.Println("No more input available.")
			break
		}
		result, err := parseAndValidateFloat(scanner.Text(), minValue, maxValue)
		if err != nil {
			fmt.Printf("%v. Please try again.\n", err)
			continue
		}
		return result
	}
	return 0
}

func parseAndValidateFloat(input string, minValue, maxValue float64) (float64, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, fmt.Errorf("input must not be empty")
	}

	result, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("value must be a valid number, got %q", trimmed)
	}

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, fmt.Errorf("value must be a finite number")
	}

	if result < minValue {
		return 0, fmt.Errorf("value must be at least %g", minValue)
	}

	if result > maxValue {
		return 0, fmt.Errorf("value must be at most %g", maxValue)
	}

	return result, nil
}
