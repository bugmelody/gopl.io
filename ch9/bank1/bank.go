// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

/**
Here’s the bank example rewritten with the balance variable confined to a monitor goroutine called teller:
 */
/**
teller ['telə] n.
1.讲述者，叙述者，讲故事者
2.记数者；计票员，点票员
3.[美国英语](银行等的)出纳员

bank tellern. 银行行员
automatic teller自动出纳机
fortune teller算命先生；预言家
automatic teller machine自动取款机；自动柜员机
automated teller machine自动柜员机
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
