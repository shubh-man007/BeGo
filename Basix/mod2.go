// Flow control statements: for, if, else, switch and defer
package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}

	return fmt.Sprint(math.Sqrt(x))
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func pow1(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	} else {
		fmt.Printf("%g >= %g\n", v, lim)
	}
	return lim
}

// func sqrt1(x float64) {
// 	epsilon := 0.0000001
// 	var z float64 = 1
// out:   //a loop can be labelled and can be used to break out of the loop
// 	for i := 0; i < 10; i++ {
// 		z -= ((z*z - x) / (2 * z))
// 		if a := math.Sqrt(x); math.Abs(a-z) <= epsilon {
// 			break out
// 		}
// 		fmt.Println(z)
// 	}
// }

func sqrt1(x float64) {
	epsilon := 0.0000001
	z := 1.0

	for {
		prevZ := z
		z -= (z*z - x) / (2 * z)
		fmt.Println(z)

		if math.Abs(z-prevZ) < epsilon {
			break
		}
	}
}

func main() {
	//for loop in Go
	sum_for := 0
	for i := 0; i <= 10; i++ {
		sum_for += i
	}
	// fmt.Println(sum_for)

	//analogous to while in C
	sum_while := 0
	for sum_while < 1000 {
		sum_while += 1
	}
	// fmt.Println(sum_while)

	//Infinite loop
	// for {

	// }

	//If statements
	// fmt.Println(sqrt(2), sqrt(-4))
	// fmt.Println(pow(3, 2, 10), pow(3, 3, 20))
	// fmt.Println(pow1(3, 2, 10), pow1(3, 3, 20))

	// for i := 1; i <= 10; i++ {
	// 	fmt.Println("---------------", i, "----------------")
	// 	sqrt1(float64(i))
	// }

	//Switch in Go
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X")
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("%s.\n", os)
	}

	fmt.Println("When is Saturday ?")
	today := time.Now().Weekday() //Returns today's day
	switch time.Saturday {
	case today + 0:
		fmt.Println("Today")
	case today + 1:
		fmt.Println("Tomorrow")
	case today + 2:
		fmt.Println("In two days")
	case today + 3:
		fmt.Println("Too far away")
	}

	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}

	// Defer (A stack is maintained and deferred items are pushed into the stack)
	// A defer statement defers the execution of a function until the surrounding function returns.
	// defer fmt.Println("World")
	// fmt.Println("hello")

	fmt.Println("Counting")
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("Done")
}
