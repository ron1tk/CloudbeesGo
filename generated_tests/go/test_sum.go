import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../..")))

package main

import (
    "testing"
)

// TestSum_ValidInputs tests the Sum function with valid inputs.
func TestSum_ValidInputs(t *testing.T) {
    t.Run("Adding positive numbers", func(t *testing.T) {
        result := Sum(2, 3)
        expected := 5
        if result != expected {
            t.Errorf("Sum(2, 3) = %d; want %d", result, expected)
        }
    })

    t.Run("Adding negative numbers", func(t *testing.T) {
        result := Sum(-2, -3)
        expected := -5
        if result != expected {
            t.Errorf("Sum(-2, -3) = %d; want %d", result, expected)
        }
    })

    t.Run("Adding positive and negative number", func(t *testing.T) {
        result := Sum(-2, 3)
        expected := 1
        if result != expected {
            t.Errorf("Sum(-2, 3) = %d; want %d", result, expected)
        }
    })

    t.Run("Adding zeros", func(t *testing.T) {
        result := Sum(0, 0)
        expected := 0
        if result != expected {
            t.Errorf("Sum(0, 0) = %d; want %d", result, expected)
        }
    })
}

// TestMinus_ValidInputs tests the Minus function with valid inputs.
func TestMinus_ValidInputs(t *testing.T) {
    t.Run("Subtracting positive numbers", func(t *testing.T) {
        result := Minus(5, 3)
        expected := 2
        if result != expected {
            t.Errorf("Minus(5, 3) = %d; want %d", result, expected)
        }
    })

    t.Run("Subtracting negative numbers", func(t *testing.T) {
        result := Minus(-3, -2)
        expected := -1
        if result != expected {
            t.Errorf("Minus(-3, -2) = %d; want %d", result, expected)
        }
    })

    t.Run("Subtracting with negative result", func(t *testing.T) {
        result := Minus(2, 5)
        expected := -3
        if result != expected {
            t.Errorf("Minus(2, 5) = %d; want %d", result, expected)
        }
    })

    t.Run("Subtracting from zero", func(t *testing.T) {
        result := Minus(0, 5)
        expected := -5
        if result != expected {
            t.Errorf("Minus(0, 5) = %d; want %d", result, expected)
        }
    })

    t.Run("Subtracting zero", func(t *testing.T) {
        result := Minus(5, 0)
        expected := 5
        if result != expected {
            t.Errorf("Minus(5, 0) = %d; want %d", result, expected)
        }
    })
}

// This set of tests ensures that for both Sum and Minus functions, all significant paths are covered.
// There's no setup or teardown needed for these tests, and there are no external dependencies to mock.