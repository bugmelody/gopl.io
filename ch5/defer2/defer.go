// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 151.

/**
Readers familiar with exceptions in other languages may be surprised that runtime.Stack
can print information about functions that seem to have already been "unwound". Go's panic
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

	// Stack formats a stack trace of the calling goroutine into buf
	// and returns the number of bytes written to buf.
	// If all is true, Stack formats stack traces of all other goroutines
	// into buf after the trace for the current goroutine.
	// ### func Stack(buf []byte, all bool) int {
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
原始结果,由于增加了注释,结果已经不适用
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


/**
运行:
$ go run ch5/defer2/defer.go 2> /dev/null
f(3)
f(2)
f(1)
defer 1
defer 2
defer 3
goroutine 1 [running]:
main.printStack()
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:47 +0x75
panic(0x49caa0, 0xc042004090)
        D:/Go/src/runtime/panic.go:458 +0x251
main.f(0x0)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:55 +0x1cf
main.f(0x1)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.f(0x2)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.f(0x3)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.main()
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:32 +0x4d




$ go run ch5/defer2/defer.go 1> /dev/null
panic: runtime error: integer divide by zero

goroutine 1 [running]:
panic(0x49caa0, 0xc042004090)
        D:/Go/src/runtime/panic.go:500 +0x1af
main.f(0x0)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:55 +0x1cf
main.f(0x1)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.f(0x2)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.f(0x3)
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:57 +0x196
main.main()
        F:/qcpj/gopl/src/gopl.io/ch5/defer2/defer.go:32 +0x4d
exit status 2


 */