// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 240.

// Crawl1 crawls web links starting with the command-line arguments.
//
// This version quickly exhausts available file descriptors
// due to excessive concurrent calls to links.Extract.
//
// Also, it never terminates because the worklist is never closed.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-crawl

//!+main
func main() {
	// worklist 记录了需要处理的 item,每个 item 是一个 []string,代表一批需要抓取的链接
	worklist := make(chan []string)

	// Start with the command-line arguments.
	// 首先从命令行参数开始抓,因此启动一个 goroutine 将 os.Args[1:] 放到 wroklist
	go func() { worklist <- os.Args[1:] }()
	/**
	想想,为什么要单独启动一个 goroutine 来初始化要抓取的链接????
	因为worklist是一个无缓冲的chan,如果没有人接收,发送时会阻塞,解决方案之一就是启动一个新的goroutine专门负责send.
	另外也可以通过有缓冲的chan来解决
	 */

	// Crawl the web concurrently.
	
	/**
	这个程序永远不会结束,即使已经从初始url找到了所有的链接.
	为了让程序结束,我们需要当worklist为空并且没有goroutine活跃的时候跳出循环
	参考 \ch8\crawl2\findlinks.go 中的做法
	 */
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			// 注意到对 seen 的操作,全部是在 main goroutine 中
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

/**
Notice that the crawl goroutine takes link as an explicit parameter to avoid the problem of
loop variable capture we first saw in Section 5.6.1. Also notice that the initial send of the command-line
arguments to the worklist must run in its own goroutine to avoid dead lock, a stuck
situation in which both the main goroutine and a crawler goroutine attempt to send to each
other while neither is receiving . An alternative solution would be to use a buffered channel.
 */

//!-main

/*
//!+output
$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
http://gopl.io/
https://golang.org/help/

https://golang.org/doc/
https://golang.org/blog/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket:
                                                        too many open files
...
//!-output
*/

/**
这样无限制的使用goroutine进行抓取很快会造成 too many open files 错误, no such host 错误也是由此引发.
 */
