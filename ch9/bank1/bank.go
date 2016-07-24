// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

/**
Here’s the bank example rewritten with the balance variable confined to a monitor goroutine called teller:
 */

// Package bank provides a concurrency-safe bank with one account.
package bank

// 这两个都是无缓冲的
var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	// balance代表用户余额,仅仅在 monitor goroutine中被使用
	var balance int // balance is confined to teller goroutine

	// 首先,这是一个死循环
	for {
		select {
		case amount := <-deposits:
		// 如果有人发送(Deposit被调用),则这个case被运行.
			balance += amount
		case balances <- balance:
			// 如果有人接受(Balance被调用),则这个case被运行. 因为balances是无缓冲的chan
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
