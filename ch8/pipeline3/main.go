// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 231.

// Pipeline3 demonstrates a finite 3-stage pipeline
// with range, close, and unidirectional channel types.
package main

import "fmt"

//!+
func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
	for v := range in {
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

//!-

/**
The call counter(naturals) implicitly converts naturals, a value of type chan int, to the type of the parameter, chan<- int. The printer(squares) call does a similar implicit conversion to <-chan int. Conversions from bidirectional to unidirectional channel types are permitted in any assignment. There is no going back, however: once you have a value of a unidirectional type such as chan<- int, there is no way to obtain from it a value of type chan int that refers to the same channel data structure.
*/
