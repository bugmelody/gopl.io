// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 130.

// The wait program waits for an HTTP server to start responding.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/**
back off
1.[美国英语]后退，退回(一段距离)：
They request that you back off 100 meters to the south.
他们要求你向南后退100米。
2.(车)略向后移动，倒车，(车)退开：
Do you mind backing off a bit, Sir, while we clear this wreck out of the way?
先生，您能把车子往后移动一点吗?我们要将这辆坏车从路上清除。
3.(因恐惧等)向后退缩，后退，退让：
The crowd backed off in terror as the police charged.
警察冲上来时，人群因害怕而后退。
4.[美国俚语]放弃(要求、立场、主张、观点、辩论等)，让步，退却，退让，屈服[亦作 back down]:
Now the time for action has arrived,it is too late to back off.
行动的时刻已到，打退堂鼓已经来不及了。
to back off from one's demands
放弃原来的要求
5.[美国俚语]停止打扰(某人)；中断(讨论、辩论等)：
She backed off from the debate when she saw her classmates were getting angry.
她见她的同班同学要生气了，就中断了辩论。
6.[美国俚语]放慢速度：
Hey, back off a little.
喂，讲慢点。
7.[美国俚语]被撵出(酒吧等公共场所)，把…撵出：
He was backed off the bar.
他被撵出了酒吧间。
 */


//!+
// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential back-off.
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	// time.Minute 是 time.Duration 类型的常量,代表 一分钟的 Duration
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

//!-

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: wait url\n")
		os.Exit(1)
	}
	url := os.Args[1]
	//!+main
	// (In function main.)
	if err := WaitForServer(url); err != nil {
		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		os.Exit(1)
	}
	//!-main
}
