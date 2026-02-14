package main

import (
	"bufio"
	"math"
	"strings"
	"testing"
)

func TestParseAndValidateFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		min      float64
		max      float64
		expected float64
		wantErr  bool
	}{
		{
			name:     "valid integer",
			input:    "42",
			min:      0,
			max:      100,
			expected: 42,
		},
		{
			name:     "valid float",
			input:    "3.14",
			min:      0,
			max:      10,
			expected: 3.14,
		},
		{
			name:     "valid negative",
			input:    "-45.5",
			min:      -90,
			max:      90,
			expected: -45.5,
		},
		{
			name:     "at minimum boundary",
			input:    "-90",
			min:      -90,
			max:      90,
			expected: -90,
		},
		{
			name:     "at maximum boundary",
			input:    "90",
			min:      -90,
			max:      90,
			expected: 90,
		},
		{
			name:     "with leading and trailing whitespace",
			input:    "  51.5  ",
			min:      -90,
			max:      90,
			expected: 51.5,
		},
		{
			name:    "empty string",
			input:   "",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "non-numeric string",
			input:   "hello",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "below minimum",
			input:   "-91",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "above maximum",
			input:   "91",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "positive infinity",
			input:   "Inf",
			min:     0,
			max:     40075,
			wantErr: true,
		},
		{
			name:    "negative infinity",
			input:   "-Inf",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:    "NaN",
			input:   "NaN",
			min:     -90,
			max:     90,
			wantErr: true,
		},
		{
			name:     "zero",
			input:    "0",
			min:      -90,
			max:      90,
			expected: 0,
		},
		{
			name:     "scientific notation",
			input:    "1e2",
			min:      0,
			max:      200,
			expected: 100,
		},
		{
			name:    "scientific notation out of range",
			input:   "1e10",
			min:     0,
			max:     40075,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndValidateFloat(tt.input, tt.min, tt.max)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got result %f", tt.input, result)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestReadFloatUntilValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		variable string
		min      float64
		max      float64
		expected float64
	}{
		{
			name:     "valid on first try",
			input:    "42.5\n",
			variable: "latitude",
			min:      -90,
			max:      90,
			expected: 42.5,
		},
		{
			name:     "invalid then valid",
			input:    "abc\n42.5\n",
			variable: "latitude",
			min:      -90,
			max:      90,
			expected: 42.5,
		},
		{
			name:     "out of range then valid",
			input:    "100\n-100\n45\n",
			variable: "latitude",
			min:      -90,
			max:      90,
			expected: 45,
		},
		{
			name:     "multiple invalid then valid",
			input:    "abc\n\nInf\n200\n-200\n50\n",
			variable: "latitude",
			min:      -90,
			max:      90,
			expected: 50,
		},
		{
			name:     "negative value accepted",
			input:    "-74.006\n",
			variable: "longitude",
			min:      -180,
			max:      180,
			expected: -74.006,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			result := readFloatUntilValid(scanner, tt.variable, tt.min, tt.max)
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("got %f, want %f", result, tt.expected)
			}
		})
	}
}
