// Structs and Pointers

package main

import "fmt"

// Struct (collection of different data types) and Arrays (collection of similar data types) are the two types of aggregates.

// Struct
type Vertex struct {
	X int
	Y int
}

var (
	v1 = Vertex{1, 2}
	v2 = Vertex{X: 1}
	v3 = Vertex{}
	v4 = &Vertex{1, 2}
)

func main() {
	var i, j int = 10, 22

	p := &i
	q := &j

	fmt.Println("i:", i)
	fmt.Println("*p:", *p)

	*p = 11

	fmt.Println("New value at p:", *p)
	fmt.Println("*q:", *q)
	*q = *q / *p
	fmt.Println("New value at q:", *q)

	v := Vertex{1, 2}
	fmt.Println(v.X, v.Y)
	a := &v
	a.X = 1e6
	fmt.Println(v.X, v.Y)

	fmt.Println(v1, v2, v3, v4)
}
