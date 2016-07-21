// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 229.

// Pipeline2 demonstrates a finite 3-stage pipeline.
package main

import "fmt"

//!+
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		// In a more complex program, it might make sense for the counter and squarer functions to defer the calls to close at the outset
		close(naturals)
	}()

	// Squarer
	go func() {
		// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
		for x := range naturals {
			squares <- x * x
		}
		// In a more complex program, it might make sense for the counter and squarer functions to defer the calls to close at the outset
		close(squares)
	}()

	// Printer (in main goroutine)
	// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
	for x := range squares {
		fmt.Println(x)
	}
}

//!-
