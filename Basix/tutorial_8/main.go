// Methods

package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

type MyFloat float64

// Go does not have classes. However, you can define methods on types.
// A method is a function with a special receiver argument.
func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Methods with pointer receivers can modify the value to which the receiver points (as Scale does here).
// Since methods often need to modify their receiver, pointer receivers are more common than value receivers.
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// You can also define methods on other types (like type aliases).
func (f MyFloat) Abs() MyFloat {
	// This method returns the absolute value of a MyFloat.
	if f < 0 {
		return -f
	}
	return f
}

func main() {
	v := Vertex{3, 4}

	// For the statement v.Abs(), the receiver is a value,
	// but Go automatically takes its address to match the pointer receiver method.
	fmt.Println(v.Abs())

	// For the statement v.Scale(10), even though v is a value and not a pointer,
	// Go interprets this as (&v).Scale(10) because Scale has a pointer receiver.
	v.Scale(10)

	fmt.Println(v)

	var f MyFloat = -10
	// This will call the Abs() method defined on MyFloat
	fmt.Println(f.Abs())
}

// When to use method with pointer recv or value recv :
// 1. If there is some sort of mutation present use a pointer recv.
// 2. Memory based optimizations can be made using pointer recv, as the instances are not copied.
// 3. If one method uses a particular type of recv, use the same recv for the rest.
