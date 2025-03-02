// Packages, variables, and functions.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
)

// How to define function in Go: (The int before '{' is used to specify what data type the function is going to return)
func add(x int, y int) int {
	return x + y
}

// Functions can return multiple values as well.
func swap(x string, y string) (string, string) {
	return y, x
}

// x and y are implicitly defined, i.e they are initialised as zero. Hence we dont need to assign them.
func split(sum int) (x int, y int) {
	x = sum * 4 / 9
	y = sum - x
	return x, y
}

var C, Java, Python bool //Because these are un-initialized they will be in their zero state, i.e is false. (C, Java, Python are all bools)

// i := 7, i cannot be assigned like this outside of a function.

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

const Pi = 3.14

// This main function is the entry point for the compiler, it would contain other function calls.
func main() {
	fmt.Println("Hello World")
	fmt.Println("My favorite number is: ", rand.Intn(10))
	fmt.Println("Pi is: ", math.Pi) //Capital P is used for Pi this is because of encapsulation. imports with lowecase are for private stuff.

	fmt.Println("The add function returns", add(4, 3), "for add(4,3).")

	// a, b := swap("Hello", "World")
	// fmt.Println("The swapped strings are", a, b)

	// x, y := split(43)
	// fmt.Println("split", x, y)

	// var i int // Because this is un-initialized this will be in its zero state, i.e is 0
	// c := 10   // Assigned c the value 10, the Go runtime decides data type using initialiser value and type cannot be changed
	// fmt.Println(i, c, C, Java, Python)

	// Printf is used because we are using formatting literals
	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	// var i int
	// var f float64
	// var b bool
	// var s string
	// fmt.Printf("%v %v %v %q\n", i, f, b, s) // %q explicity prints the quotes

	// We need to explicitly perform type conversions.
	var x, y int = 3, 4
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z)

	// We can use constants to define hardcoded values in our code.
	const World = "世界"
	const Truth = true
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")
	fmt.Println("Universal", Truth)
}

// Notes: (Refer : https://go.dev/tour/list)
// 1. The go.mod file is like a requirements.txt file helps others n about the imports and their versions.
// 2. Go uses UTF-8 encoding and supports different languages literal/characters as well.
// 3. Go generates a "Cross Platform Single Binary Executable" file. (use : go build main.go command), which helps the code to run on different OS and architectures.
// 4. We use go run command to run because it uses a temporary file and runs our code without generating any executable file, comparatively faster than the build executable.
// 5. Go is a statically typed language, implying all types are verified and checked during compile time.
// 6. In Go a variable is never in an un-initialized state (its always initialized to its data types's zero-state by default)
// 7. Constants value cannot be changed, they are evaluated during compile time. Constants can be string, boolean or number.
