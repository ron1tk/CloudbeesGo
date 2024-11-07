import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

package main

import (
    "testing"
)

// TestSum_PositiveNumbers tests the Sum function with positive numbers.
func TestSum_PositiveNumbers(t *testing.T) {
    result := Sum(5, 3)
    expected := 8
    if result != expected {
        t.Errorf("Sum(5, 3) = %d; want %d", result, expected)
    }
}

// TestSum_NegativeNumbers tests the Sum function with negative numbers.
func TestSum_NegativeNumbers(t *testing.T) {
    result := Sum(-5, -3)
    expected := -8
    if result != expected {
        t.Errorf("Sum(-5, -3) = %d; want %d", result, expected)
    }
}

// TestSum_PositiveAndNegativeNumbers tests the Sum function with both positive and negative numbers.
func TestSum_PositiveAndNegativeNumbers(t *testing.T) {
    result := Sum(5, -3)
    expected := 2
    if result != expected {
        t.Errorf("Sum(5, -3) = %d; want %d", result, expected)
    }
}

// TestSum_Zero tests the Sum function with zero.
func TestSum_Zero(t *testing.T) {
    result := Sum(0, 0)
    expected := 0
    if result != expected {
        t.Errorf("Sum(0, 0) = %d; want %d", result, expected)
    }
}

// TestMinus_PositiveNumbers tests the Minus function with positive numbers.
func TestMinus_PositiveNumbers(t *testing.T) {
    result := Minus(5, 3)
    expected := 2
    if result != expected {
        t.Errorf("Minus(5, 3) = %d; want %d", result, expected)
    }
}

// TestMinus_NegativeNumbers tests the Minus function with negative numbers.
func TestMinus_NegativeNumbers(t *testing.T) {
    result := Minus(-5, -3)
    expected := -2
    if result != expected {
        t.Errorf("Minus(-5, -3) = %d; want %d", result, expected)
    }
}

// TestMinus_PositiveAndNegativeNumbers tests the Minus function with both positive and negative numbers.
func TestMinus_PositiveAndNegativeNumbers(t *testing.T) {
    result := Minus(5, -3)
    expected := 8
    if result != expected {
        t.Errorf("Minus(5, -3) = %d; want %d", result, expected)
    }
}

// TestMinus_Zero tests the Minus function with zero.
func TestMinus_Zero(t *testing.T) {
    result := Minus(0, 0)
    expected := 0
    if result != expected {
        t.Errorf("Minus(0, 0) = %d; want %d", result, expected)
    }
}

// TestMinus_LargeNumbers tests the Minus function with large numbers.
func TestMinus_LargeNumbers(t *testing.T) {
    result := Minus(1000000, 500000)
    expected := 500000
    if result != expected {
        t.Errorf("Minus(1000000, 500000) = %d; want %d", result, expected)
    }
}