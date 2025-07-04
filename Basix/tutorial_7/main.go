// Maps and function closures

package main

import (
	"fmt"
	"strings"
)

type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

func WordCount(s string) map[string]int {
	var str_list []string = strings.Fields(s)
	m := make(map[string]int)
	for _, v := range str_list {
		m[v] += 1
	}
	return m
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
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{40.88, -74.39}
	fmt.Println(m["Bell Labs"])

	// map literal
	m1 := map[string]Vertex{
		"a": {1.1, 2.2},
		"b": {3.3, 4.4},
	}
	fmt.Println(m1)

	m2 := make(map[string]string)
	m2["Answer"] = "Yes"
	fmt.Println(m2["Answer"])

	m2["Answer"] = "No"
	fmt.Println(m2["Answer"])

	m2["Hello"] = "World"

	delete(m2, "Answer")
	fmt.Println(m2)

	v, ok := m2["Hello"]
	fmt.Println(v, ok)

	fmt.Println(WordCount("The dog ate the cookies"))

	p := adder()
	fmt.Println(p(10))

	q := fibonacci()
	for i := 0; i < 5; i++ {
		fmt.Println(q())
	}
}
