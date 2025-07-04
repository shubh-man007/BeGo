// Slices
// Slices are a reference type because they have a backing array they refer to.
// Zero value for a reference type is null

package main

import (
	"fmt"
)

func printSlice(s []int) {
	fmt.Printf("len = %d, cap = %d, %v\n", len(s), cap(s), s)
}

func stringSlice(x string, s []int) {
	fmt.Printf("%v, len = %d, cap = %d, %v\n", x, len(s), cap(s), s)
}

func main() {
	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	var s1 []int = primes[:2]
	var s2 []int = primes[1:5]

	fmt.Println(s1, s2)

	// Because slices are reference types, changing an element would effect the index in the backing array as well.
	s1[1] = 17
	fmt.Println(s1, s2)
	fmt.Println(primes)

	// Ways of defining a struct:
	// Slice literal:
	q := []int{1, 2, 3}
	fmt.Println(q)

	// Slice of structs:
	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s)

	// Capacity is defined as length from starting point to the end of the backing array
	// Length is the number of blocks between starting point and ending point.
	s_q := primes[:0]
	printSlice(s_q)

	s_q = s_q[1:4]
	printSlice(s_q)

	s_q = s_q[:2]
	printSlice(s_q)

	s_q = s_q[1:]
	printSlice(s_q)

	// // zero value of a slice is nil.
	// var x []int
	// fmt.Println(x, len(x), cap(x))
	// if x == nil {
	// 	fmt.Println("nil!")
	// }

	// make function to initialize slices, channels, maps, etc.:
	a := make([]int, 5)
	stringSlice("a", a)

	b := make([]int, 0, 5)
	stringSlice("b", b)

	c := b[:3]
	stringSlice("c", c)

	d := c[1:3]
	stringSlice("d", d)

	e := d[1:4]
	stringSlice("e", e)

	// we can append one or mre than one value to an existing slice. append is a variadic function.
	f := []int{}
	stringSlice("f", f)

	f = append(f, 0, 1, 2, 3)
	stringSlice("f", f)

	f = append(f, 5)
	stringSlice("f", f)

	// If we know before hand the capacity of a slice, we can use the "make" function to initialize the slice and then append to it.
	// This helps reduce computations for append by pre determining the slice capacity and append doesnt have to search and double the capacity (threshold 256) every time len == cap.

	g := make([]int, 0, 5)
	for i := 0; i < 5; i++ {
		g = append(g, i)
		printSlice(g)
	}

	// The range form of the for loop iterates over a slice or map.
	for i, v := range primes {
		fmt.Printf("Index : %v, Value : %v\n", i, v)
	}
}
