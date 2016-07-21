// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		// reverb1 中,每个连接是单独的goroutine处理,但是每个echo只能顺序输出
		// 因此在reverb2中,对每一个client的输入,单独使用goroutine进行echo处理
		/**
		注意参考 net.Conn 的文档,里面明确说明了多个goroutine可以同时调用同一个 Conn 的方法
		
		// Multiple goroutines may invoke methods on a Conn simultaneously.
		type Conn interface {
		 */
		
		/**
		The arguments to the function started by go are evaluated when the go statement itself is executed;
		thus input.Text() is evaluated in the main goroutine.
		 */
		go echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
