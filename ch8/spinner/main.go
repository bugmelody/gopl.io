// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 218.

// Spinner displays an animation while computing the 45th Fibonacci number.
package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			// %c: the character represented by the corresponding Unicode code point
			//  \r 是什么,用于实现console中字符转圈的效果
			fmt.Printf("\r%c", r)

			/**
			// Sleep pauses the current goroutine for at least the duration d.
			// A negative or zero duration causes Sleep to return immediately.
			func Sleep(d Duration)
			*/
			time.Sleep(delay)
		}
	}
}

/**
斐波那契数列指的是这样一个数列 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55
这个数列从第3项开始，每一项都等于前两项之和。
*/
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

//!-

/**
The main function then returns. When this happens, all goroutines are abruptly terminated and the program exits. Other than by returning from main or exiting the program, there is no programmatic way for one goroutine to stop another, but as we will see later, there are ways to communicate with a goroutine to request that it stop itself.

Notice how the program is expressed as the composition of two autonomous activities, spinning and Fibonacci computation. Each is written as a separate function but both make progress concurrently. 
 */