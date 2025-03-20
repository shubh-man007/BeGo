// Methods

package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X float64
	Y float64
}

type MyFloat float64

// Go does not have classes. However, you can define methods on types.
// A method is a function with a special receiver argument.
func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Methods with pointer receivers can modify the value to which the receiver points (as Scale does here).
// Since methods often need to modify their receiver, pointer receivers are more common than value receivers.
func (v *Vertex) Scale(f float64) {
	v.X = f * v.X
	v.Y = f * v.Y
}

func (m MyFloat) Abs() float64 {
	if m < 0 {
		return float64(-m)
	}
	return float64(m)
}

func main() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
	// For the statement v.Scale(5), even though v is a value and not a pointer, the method with the pointer receiver is called automatically.
	// Go interprets the statement v.Scale(5) as (&v).Scale(5) since the Scale method has a pointer receiver.
	// Methods with value receivers take either a value or a pointer as the receiver when they are called
	v.Scale(10)
	fmt.Println(v.Abs())
}
