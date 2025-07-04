// Interfaces

package main

import (
	"fmt"
	"math"
)

// a MyFloat implements Abser
// a *Vertex implements Abser
type Abser interface {
	Abs() float64
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

type I interface {
	M()
}

type T struct {
	S string
}

func (t T) M() {
	fmt.Println(t.S)
}

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice of %v is %v\n", v, v*2)

	case string:
		fmt.Printf("%v is %v bytes long\n", v, len(v))

	default:
		fmt.Println("No clue")
	}

}

func main() {
	var a Abser

	v := Vertex{3, 4}
	a = &v
	fmt.Println(a.Abs())
	// a = v (would not work as vertex's Abs() method takes in a pointer recv.)

	f := MyFloat(10)
	a = f
	fmt.Println(a)

	// Here, i now has now been assigned a dynamic type T and a dynamic value "hello"
	var i I = T{"hello"}
	i.M()

	// Empty Interface:
	var p interface{}
	describe(p)

	p = 42
	describe(p)

	p = "hello"
	describe(p)

	// Type assertion
	p = "World"

	s := p.(string)
	fmt.Println(s)

	s, ok := p.(string)
	fmt.Println(s, ok)

	r, ok := p.(float64)
	fmt.Println(r, ok)

	do(21)
	do("Hello")
	do(true)

}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
