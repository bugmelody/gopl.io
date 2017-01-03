// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 243.

// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!+
func main() {
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() { worklist <- os.Args[1:] }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				/**
				Links found by crawl are sent to the worklist from a dedicated goroutine to avoid deadlock.

				这里为什么新抓到的链接要单独用一个goroutine来发送呢?
				因为 worklist 是无缓冲的 channel,如果多个goroutine同时send,而在main goroutine 中接收太慢,就会引起一段时间的阻塞,导致
				worklist <- foundLinks 会停顿一段时间,直到被接收

				总之, 针对无缓冲的channel,如果要向它send数据,很有可能会被阻塞(永久或暂时),如果想避免阻塞,要么单独起
				一个goroutine进行send,要么使用有缓冲的channel
				*/
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	// seen这个map被限制在只从main goroutine中使用
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

//!-

/**
The crawler goroutines are all fed by the same channel, unseenLinks. The main goroutine is
responsible for de-duplicating items it receives from the worklist, and then sending each
unseen one over the unseenLinks channel to a crawler goroutine.

The seen map is confined within the main goroutine; that is, it can be accessed only by that
goroutine. Like other forms of information hiding, confinement helps us reason about the
correctness of a program. For example, local variables cannot be mentioned by name from
outside the function in which they are declared; variables that do not escape (§2.3.4) from a
function cannot be accessed from outside that function; and encapsulated fields of an object
cannot be accessed except by the methods of that object. In all cases, information hiding helps
to limit unintended interactions between parts of the program.

Links found by crawl are sent to the worklist from a dedicated goroutine to avoid dead lock.
To save space, we have not addressed the problem of termination in this example.
 */