package main

import (
	"testing"
)

// TestSum_PositiveNumbers tests the Sum function with positive numbers.
func TestSum_PositiveNumbers(t *testing.T) {
	result := Sum(2, 3)
	expected := 5
	if result != expected {
		t.Errorf("Sum(2, 3) = %d; want %d", result, expected)
	}
}

// TestSum_NegativeNumbers tests the Sum function with negative numbers.
func TestSum_NegativeNumbers(t *testing.T) {
	result := Sum(-2, -3)
	expected := -5
	if result != expected {
		t.Errorf("Sum(-2, -3) = %d; want %d", result, expected)
	}
}

// TestSum_Zero tests the Sum function with zero values.
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

// TestMinus_Zero tests the Minus function when the result is zero.
func TestMinus_Zero(t *testing.T) {
	result := Minus(2, 2)
	expected := 0
	if result != expected {
		t.Errorf("Minus(2, 2) = %d; want %d", result, expected)
	}
}

// TestMinus_ResultNegative tests the Minus function when the result is negative.
func TestMinus_ResultNegative(t *testing.T) {
	result := Minus(2, 5)
	expected := -3
	if result != expected {
		t.Errorf("Minus(2, 5) = %d; want %d", result, expected)
	}
}