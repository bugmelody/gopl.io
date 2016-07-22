// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		// 一旦从stdin读取到任何数据,就向about发送事件
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	/**
	We can’t just receive from each channel because whichever operation we try first will block until
	completion. We need to multiplex these operations, and to do that, we need a select statement.
	如果采用普通的receive操作,不论先从tick或者是about channel中接收,都会造成阻塞,因此需要使用 select.
	*/
	/**
	A select waits until a communication for some case is ready to proceed. It then performs
	that communication and executes the case’s associated statements; the other communications
	do not happen. A select with no cases, select{}, waits forever.
	*/
	
	/**
	 The time.After function immediately returns a channel, and starts a new goroutine that sends a single value on that channel after the specified time. The select statement below waits until the first of two events arrives, either an abort event or the event indicating that 10 seconds have elapsed. If 10 seconds go by with no abort, the launch proceeds. 
	 */
	select {
	// time.After 会立即返回一个chan,然后启动一个新的goroutine在指定的时间后发送数据到返回的这个chan
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		// 如果接收到about事件,从main返回,launch被取消
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
