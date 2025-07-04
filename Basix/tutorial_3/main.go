// Loops, Conditionals, Defer, Switch

package main

import (
	"errors"
	"fmt"
	"math"
)

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}

	return fmt.Sprint(math.Sqrt(x)) //Sprint returns the formatted string
}

func powlim(base, power, lim float64) (float64, error) {
	var err error
	if v := math.Pow(base, power); v < lim {
		return v, err
	} else {
		err := errors.New("exceeded computational limit")
		return lim, err
	}
}

func sqrtest(x float64) float64 {
	var z float64 = 1
	var epsilon float64 = 0.0001
	for math.Abs(z*z-x) > epsilon {
		fmt.Println(z)
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func func1(x int64) int64 {
	var i int64
	for i = 1; i < 100000; i++ {
		x += i
	}
	return x
}

func main() {

	defer fmt.Println(func1(0))

	var sum1 int = 0
	for i := 0; i < 10; i++ {
		sum1 += i
	}
	fmt.Printf("The sum is %v\n", sum1)

	var sum2 int = 1
	for sum2 < 1000 {
		sum2 += sum2
	}
	fmt.Printf("The sum is %v\n", sum2)

	fmt.Println(sqrt(10))

	ans, err := powlim(10, 3, 10000)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ans)

	fmt.Println(sqrtest(10))

	var guess int = 11
	var actual int = 13

	switch actual {
	case guess + 1:
		fmt.Println("Just one away")
	case guess + 2:
		fmt.Println("Skip away")
	default:
		fmt.Println("Nevermind")
	}

	for i := 0; i < 100; i++ {
		if i == 99 {
			fmt.Println("Loop Done")
		}
	}
}
