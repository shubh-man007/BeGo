//Interfaces
// An interface type is defined as a set of method signatures.
// A value of interface type can hold any value that implements those methods.

package main

import (
	"fmt"
	"math"
	"time"
)

type I interface {
	M()
}

type MyType struct {
	S string
}

// This method means type MyType implements the interface I,
// but we don't need to explicitly declare that it does so.
func (m *MyType) M() {
	if m == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(m.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%q is %v bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
}

type Myerror struct {
	When time.Time
	What string
}

func (m *Myerror) Error() string {
	return fmt.Sprintf("at %s (%v)", m.What, m.When)
}

func run() error {
	return &Myerror{time.Now(), "An error occured"}
}

type Myname struct {
	Name string
	Age  int
}

// Stringer
func (m Myname) String() string {
	return fmt.Sprintf("%v (%v years)", m.Name, m.Age)
}

func main() {
	var i I = &MyType{"Hello"}
	// var i I
	i.M()
	describe(i)

	// If the concrete value inside the interface itself is nil, the method will be called with a nil receiver.
	var t *MyType
	i = t
	describe(i)
	i.M()

	f := F(math.Pi)
	var j I = f
	j.M()
	describe(j)

	// An interface can be conceptualized as a container which contains a dynamic type and a dynamic value.
	// In this case the interface 'i' contains a dynamic type of MyType and a  dynamic value of "Hello".
	// Calling a method on an interface value executes the method of the same name on its underlying type.
	// The zero-state of an interface is nil. If we comment out "var i I = MyType{"Hello"}" and run "var i I" instead Go throws a panic.
	// panic: runtime error: invalid memory address or nil pointer dereference

	var i_nil I
	// i_nil.M()
	describe(i_nil)

	// The interface type that specifies zero methods is known as the empty interface. An empty interface may hold values of any type.
	var a interface{}
	describe_a(a)

	a = 10
	describe_a(a)

	a = "Hello"
	describe_a(a)

	// Type Assertion
	var b interface{} = "Hello"

	// This statement asserts that the interface value b holds the concrete type "string" and assigns the underlying "string" value ("Hello") to the variable s.
	s := b.(string)
	fmt.Println(s)

	// If b holds a "string", then s will be the underlying value and ok will be true.
	// If not, ok will be false and s will be the zero value of type "string", and no panic occurs.
	s, ok := b.(string)
	fmt.Println(s, ok)

	do(1)
	do("Hello")
	do(math.Pi)

	if err := run(); err != nil {
		fmt.Println(err)
	}

	x := Myname{"Posto", 21}
	y := Myname{"Geetu", 22}
	z := Myname{"Chello", 22}

	fmt.Println(x, y, z)

}

func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func describe_a(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
