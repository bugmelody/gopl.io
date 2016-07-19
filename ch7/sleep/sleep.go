// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 179.

// The sleep program sleeps for a specified period of time.
package main

import (
	"flag"
	"fmt"
	"time"
)

//!+sleep
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	/** The fmt package calls the time.Duration’s String method to print the period
	not as a number of nanoseconds, but in a user-friendly notation
	$ go doc time.Duration.String # 看看
	*/
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

//!-sleep
