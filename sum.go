package main

import "fmt"

// Sum adds two integers and returns the result.
func Sum(a int, b int) int {
    return a + b
}
func Minus(a int, b int) int {
    return a - b
}

func multiply(a int, b int) int {
    return a * b
}





func main() {
    fmt.Println("Sum of 5 and 3 is", Sum(5, 3))
}
