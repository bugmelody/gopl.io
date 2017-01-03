// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	fmt.Println("Commencing countdown.")
	
	// ### func Tick(d Duration) <-chan Time {
	tick := time.Tick(1 * time.Second)
	// 每隔 1s 会从 tick 这个 channel 收到时间信号
	for countdown := 10; countdown > 0; countdown-- {
		// 每隔 1s 输出倒计时
		fmt.Println(countdown)
		<-tick
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}

/*
The time.Tick function returns
a channel on which it sends events periodically, acting like a metronome. The value of each
event is a timestamp, but it is rarely as interesting as the fact of its delivery.
 */