// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 151.

/**
Readers familiar with exceptions in other languages may be surprised that runtime.Stack
can print information about functions that seem to have already been ‘‘unwound.’’ Go’s panic
mechanism runs the deferred functions before it unwinds the stack.
 */

/**
运行
查看 stdout
$ go run ch5/defer2/defer.go 2> /dev/null
查看 stderr
$ go run ch5/defer2/defer.go 1> /dev/null
 */

// Defer2 demonstrates a deferred call to runtime.Stack during a panic.
package main

import (
	"fmt"
	"os"
	"runtime"
)

//!+
func main() {
	defer printStack()
	f(3)
}

// printStack 打印出来的 stack 信息是主动打印出来的,并不是 panic 出来的,
// panic 出来的信息到了 stderr, printStack 打印的信息在 stdout
func printStack() {
	// 声明一个 byte 数组,长度为 4096
	var buf [4096]byte
	// runtime.Stack 要求第一个参数是 []byte, 这里必须进行切片操作 
	n := runtime.Stack(buf[:], false)
	// runtime.Stack 写了多少数据到 buf 是通过 runtime.Stack 的返回值来反应的
	os.Stdout.Write(buf[:n])
}

//!-

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

/*
//!+printstack
goroutine 1 [running]:
main.printStack()
	src/gopl.io/ch5/defer2/defer.go:20
main.f(0)
	src/gopl.io/ch5/defer2/defer.go:27
main.f(1)
	src/gopl.io/ch5/defer2/defer.go:29
main.f(2)
	src/gopl.io/ch5/defer2/defer.go:29
main.f(3)
	src/gopl.io/ch5/defer2/defer.go:29
main.main()
	src/gopl.io/ch5/defer2/defer.go:15
//!-printstack
*/
