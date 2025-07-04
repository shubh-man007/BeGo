// Generics

package main

import (
	"fmt"
)

// We need to repeat and modify the boiler plate code to execute similar functions of int and string.
func FindIndexInt(s []int, a int) int {
	for i, v := range s {
		if v == a {
			return i
		}
	}
	return -1
}

func FindIndexStr(s []string, a string) int {
	for i, v := range s {
		if v == a {
			return i
		}
	}
	return -1
}

// This declaration means that s is a slice of any type T that fulfills the built-in constraint comparable. a is also a value of the same type.
// Comparable is a useful constraint that makes it possible to use the == and != operators on values of the type.
func Index[T comparable](s []T, a T) int {
	for i, v := range s {
		if v == a {
			return i
		}
	}
	return -1
}

// In addition to generic functions, Go also supports generic types.
// A type can be parameterized with a type parameter, which could be useful for implementing generic data structures.
type List[T any] struct {
	next *List[T]
	val  T
}

func LLTraversal[T any](l *List[T]) {
	for l != nil {
		fmt.Println(l.val)
		l = l.next
	}
}

func NewList[T any](values ...T) *List[T] {
	if len(values) == 0 {
		return nil
	}
	head := &List[T]{val: values[0]}
	current := head
	for _, v := range values[1:] {
		current.next = &List[T]{val: v}
		current = current.next
	}
	return head
}

func main() {
	s1 := []int{1, 2, 3, 4, 5, 7, 5}
	s2 := []string{"Shubh", "Geetu", "Chello"}

	fmt.Println(FindIndexInt(s1, 5))
	fmt.Println(FindIndexStr(s2, "Bro"))

	fmt.Println(Index(s1, 5))
	fmt.Println(Index(s2, "Bro"))

	head := NewList(1, 2, 3, 4, 5)
	LLTraversal(head)
}
