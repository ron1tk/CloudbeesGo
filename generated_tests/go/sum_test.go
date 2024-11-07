package main

import (
    "testing"
)

// TestSum tests the Sum function with various inputs.
func TestSum(t *testing.T) {
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"two positive numbers", 5, 3, 8},
        {"positive and negative number", -5, 3, -2},
        {"two negative numbers", -5, -3, -8},
        {"zero and positive number", 0, 5, 5},
        {"zero and negative number", 0, -5, -5},
        {"two zeros", 0, 0, 0},
        {"large numbers", 1000000, 500000, 1500000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Sum(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Sum(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}

// TestMinus tests the Minus function with various inputs.
func TestMinus(t *testing.T) {
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"positive greater than negative", 5, 3, 2},
        {"negative greater than positive", 3, 5, -2},
        {"two negative numbers", -5, -3, -2},
        {"negative and positive number", -5, 3, -8},
        {"zero and positive number", 0, 5, -5},
        {"positive number and zero", 5, 0, 5},
        {"two zeros", 0, 0, 0},
        {"large numbers", 1000000, 500000, 500000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Minus(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Minus(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}