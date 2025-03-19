// Types of reference types in Go:
// 1. Pointers
// 2. Slices
// 3. Functions
// 4. Channel
// 5. Map
// PS : Zero value of Reference types in nil

package main

import (
	"fmt"
	"math"
	"strings"
	// "golang.org/x/tour/pic"
	// "golang.org/x/tour/wc"
)

// Structs (A collection of fields)
type Vertex struct {
	X int
	Y int
}

// var (
// 	v1    = Vertex{1, 2}
// 	v2    = Vertex{X: 1}
// 	v3    = Vertex{}
// 	p_var = &Vertex{1, 2} // has type *Vertex
// )

func printSlice(s []int) {
	fmt.Printf("len = %d, cap = %d, %v\n", len(s), cap(s), s)
}

func printSliceString(s string, x []int) {
	fmt.Printf("%s, len = %d, cap = %d, %v\n", s, len(x), cap(x), x)
}

func Pic(dx, dy int) [][]uint8 {
	slice_Pic := make([][]uint8, dy)

	for y := 0; y < dy; y++ {
		slice_Pic[y] = make([]uint8, dx)
		for x := 0; x < dx; x++ {
			slice_Pic[y][x] = uint8(math.Sin(float64(x)) * float64(y))
		}
	}

	return slice_Pic
}

type Location struct {
	Lat int
	Log int
}

var n = map[string]Location{
	"Bell Labs": Location{
		40, -74,
	},
	"Google": Location{
		37, -122,
	},
}

func WordCount(s string) map[string]int {
	words_list := strings.Fields(s)
	m := make(map[string]int)
	for _, value := range words_list {
		m[value] += 1
	}
	return m
}

func Compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}


func main() {
	// Pointers:
	i, j := 10, 20

	p := &i
	fmt.Println(*p)
	*p = 21
	fmt.Println(i)

	p = &j
	*p = *p / 2
	fmt.Println(j)

	// Structs
	v := Vertex{1, 2}
	fmt.Println(v)
	v.X = 4 // Chainging the value of field X
	fmt.Println(v.X)

	p_vertex := &v
	fmt.Println(*p_vertex)
	p_vertex.Y = 8
	fmt.Println(v)

	fmt.Println(v1, v2, v3, *p_var)

	// Arrays
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a)
	a_text := a[0] + " " + a[1]
	fmt.Println(a_text)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	// Slices
	names := [4]string{"John", "Paul", "George", "Ringo"}
	fmt.Println(names)

	a := names[0:3] // or a := names[:3]
	b := names[1:3]
	fmt.Println(a, b)

	// Slice is a reference type. Names, a and b all point to the same memory location/data.
	// Hence they all have the same value at the same indices and changing value at one index would change it for names a and b.
	b[0] = "Kevin"
	fmt.Println(a, b)
	fmt.Println(names)

	// Slice Literal
	// We can declare slice directly without declaring a backing array.
	q := []int{2, 3, 4, 5, 6, 6}
	printSlice(q)

	r := []bool{true, false, true, false, false, true}
	fmt.Println(r)

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{6, false},
		{6, false},
	}
	fmt.Println(s)

	// Capacity is defined as length from starting point to the end of the backing array
	// Length is the number of blocks between starting point and ending point.
	s_q := q[:0]
	printSlice(s_q)

	s_q = q[1:4]
	printSlice(s_q)

	s_q = q[:2]
	printSlice(s_q)

	s_q = q[1:]
	printSlice(s_q)

	// zero value of a slice is nil.
	var x []int
	fmt.Println(x, len(x), cap(x))
	if x == nil {
		fmt.Println("nil!")
	}

	// make is a function provided by Go, to init reference types.
	// make{type(int slice, etc.), length, capacity}
	//capacity >= length
	a := make([]int, 5)
	printSliceString("a", a)

	b := make([]int, 0, 5)
	printSliceString("b", b)

	c := b[:2]
	printSliceString("c", c)

	d := c[2:5]
	printSliceString("d", d)

	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}

	// Excercise-1 :
	// pic.Show(Pic)

	append_s := []int{}
	printSlice(append_s)

	// we can append one or mre than one value to an existing slice. append is a variadic function.
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			append_s = append(append_s, i)
			printSlice(append_s)
		} else {
			append_s = append(append_s, i, i+1, i+2)
			printSlice(append_s)
			i += 3
		}
	}

	// If we know before hand the capacity of a slice, we can use the "make" function to initialize the slice and then append to it.
	// This helps reduce computations for append by pre determining the slice capacity and append doesnt have to search and double the capacity (threshold 256) every time len == cap.

	// Maps (Dont use var to initialize a map)
	m := make(map[string]Location)
	m["Hello World"] = Location{10, 20}
	// fmt.Println(m["Hello World"])

	for key, value := range m {
		fmt.Println(key, value)
	}

	fmt.Println(n)

	// Excercise-2 :
	// wc.Test(WordCount)

	hypot := func(x, y float64) float64 {
		return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	}

	hypot_test := hypot(8, 6)
	fmt.Println(hypot_test)

	fmt.Println(Compute(hypot))

	pos := adder()

	for i := 0; i < 3; i++ {
		fmt.Println(pos(i))
	}

	// Excercise-3 :
	fib := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(fib())
	}

}
