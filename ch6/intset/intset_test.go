// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package intset

import "fmt"

func Example_one() {
	//!+main
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"
	//!-main

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func Example_two() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	/**
	A word of caution: we declared String and Has as methods of the pointer type *IntSet not
	out of necessity, but for consistency with the other two methods, which need a pointer receiver
	because they assign to s.words. Consequently, an IntSet value does not have a String
	method, occasionally leading to surprises like this:
	 */
	//!+note
	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"
	//!-note

	/**
	In the first case, we print an *IntSet pointer, which does have a String method. In the
	second case, we call String() on an IntSet variable; the compiler inserts the implicit & op er-
	ation, giving us a pointer, which has the String method. But in the third case, because the
	IntSet value does not have a String method, fmt.Println prints the representation of the
	struct instead. It’s important not to forget the & operator. Making String a method of
	IntSet, not *IntSet, might be a good idea, but this is a case-by-case judgment.
	*/

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
}
