// Arrays (fixed type)
// Go is a statically typed language and as size of array is part of its type, hence array cannot be resized.
// Zero value of an array is a zero array (i.e, all the values in a array in the zero state)

package main

import (
	"fmt"
)

func main() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a)

	var b [5]int
	for i := 0; i < 5; i++ {
		b[i] = i + 1
	}
	fmt.Println(b)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)
}
