// Data Types and Functions

package main

import (
	"errors"
	"fmt"
	"math/cmplx"
)

var C, Java, Python bool

const Pi = 3.14

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

func main() {

	var x string = "Hello"
	var y string = "World"

	fmt.Println(x, y)
	fmt.Println(swap(x, y))

	// Because both num1 and Pi are untyped constants, the compiler allows them to be multiplied together and assigns the result an appropriate type (float64 in this case, because of 3.14).
	// This is part of Go’s constant rules:
	// Untyped constants do not have a fixed type until they are used in a context that requires one.
	// When used in an expression, Go infers the type based on the other operand.
	const num1 = 10
	// Assigned c the value 10, the Go runtime decides data type using initialiser value and type cannot be changed
	c := 11
	fmt.Printf("Type : %T, Value : %v, num1\n", num1, num1)
	fmt.Printf("Type : %T, Value : %v, Pi\n", Pi, Pi)
	fmt.Printf("Type : %T, Value : %v, c\n", c, c)

	fmt.Println(num1 * Pi)
	// fmt.Println(c * Pi)
	fmt.Println(c * num1)

	var sum1 float64 = 18
	var a, b float64 = split(sum1)
	// fmt.Println("The split function returns : ", a, "and", b, "for split(18)")
	fmt.Printf("The split function returns: %v, %v for split(18).\n", a, b)

	var i int // Because this is un-initialized this will be in its zero state, i.e is 0
	fmt.Println(i, c, C, Java, Python)

	fmt.Println("The add function returns", add(4, 3), "for add(4,3).")

	fmt.Printf("Type : %T, Value : %v\n", ToBe, ToBe)
	fmt.Printf("Type : %T, Value : %v\n", MaxInt, MaxInt)
	fmt.Printf("Type : %T, Value : %v\n", z, z)

	var quo, rem, err = intDiv(10, 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("The quotient is: %v and the remainder is: %v", quo, rem)
}

func swap(x string, y string) (string, string) {
	return y, x
}

func split(sum float64) (x float64, y float64) {
	x = sum * 4 / 9
	y = sum - x
	return x, y
}

func add(x int, y int) int {
	return x + y
}

func intDiv(a int, b int) (int, int, error) {
	var err error
	if b == 0 {
		err = errors.New("cannot divide by zero")
		return 0, 0, err
	}
	var x int = a / b
	var y int = a % b
	return x, y, err
}

// Notes: (Refer : https://go.dev/tour/list)
// 1. The go.mod file is like a requirements.txt file helps others n about the imports and their versions.
// 2. Go uses UTF-8 encoding and supports different languages literal/characters as well.
// 3. Go generates a "Cross Platform Single Binary Executable" file. (use : go build main.go command), which helps the code to run on different OS and architectures.
// 4. We use go run command to run because it uses a temporary file and runs our code without generating any executable file, comparatively faster than the build executable.
// 5. Go is a statically typed language, implying all types are verified and checked during compile time.
// 6. In Go a variable is never in an un-initialized state (its always initialized to its data types's zero-state by default)
// 7. Constants value cannot be changed, they are evaluated during compile time. Constants can be string, boolean or number.
