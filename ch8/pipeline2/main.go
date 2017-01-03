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
		// at the outset: 开始；起初
		close(naturals)
	}()

	// Squarer
	go func() {
		// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
		for x := range naturals {
			squares <- x * x
		}
		// In a more complex program, it might make sense for the counter and squarer functions to defer the calls to close at the outset
		// at the outset: 开始；起初
		close(squares)
	}()

	// Printer (in main goroutine)
	// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
	for x := range squares {
		fmt.Println(x)
	}
}

//!-

/**
You needn't close every channel when you've finished with it. It's only necessary to close a
channel when it is important to tell the receiving goroutines that all data have been sent. A
channel that the garbage collector determines to be unreachable will have its resources
reclaimed whether or not it is closed. (Don't confuse this with the close operation for open
files. It is important to call the Close method on every file when you've finished with it.)
 */
