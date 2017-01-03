// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 246.

// Countdown implements the countdown for a rocket launch.
package main

// NOTE: the ticker goroutine never terminates if the launch is aborted.
// This is a "goroutine leak".

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		// 一旦读取到一个字节后,发送 abort 通知
		abort <- struct{}{}
	}()

	//!+
	// commence [kə'mens] vi.开始: vt.使开始；着手
	fmt.Println("Commencing countdown.  Press return to abort.")
	/**
	The time.Tick function behaves as if it creates a goroutine that calls time.Sleep in a loop, sending an
	event each time it wakes up. When the main function returns, it stops receiving events from tick, but
	the ticker goroutine is still there, trying in vain to send on a channel from which no goroutine is
	receiving — a goroutine leak.

	The Tick function is convenient, but it's appropriate only when the ticks will be needed throughout
	the lifetime of the application. Otherwise, we should use this pattern:

	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C // receive from the ticker's channel
	ticker.Stop() // cause the ticker's goroutine to terminate
	
	不再需要ticker的时候,应该调用ticker.Stop(),否则会造成goroutine leak,这里是因为整个程序退出所以使用time.Tick(xxx)更方便.
	*/
	tick := time.Tick(1 * time.Second)
	
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		// The select statement below causes each iteration of the loop to wait up to 1 second for an abort, but no longer.
		select {
		case <-tick:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	launch()
}

//!-

func launch() {
	// lift off （火箭等）发射；（直升飞机）起飞
	fmt.Println("Lift off!")
}
