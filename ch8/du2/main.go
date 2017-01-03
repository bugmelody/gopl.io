// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 249.

// The du2 command computes the disk usage of the files in a directory.
package main

// The du2 variant uses select and a time.Ticker
// to print the totals periodically if -v is set.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//!+
var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// ...start background goroutine...

	//!-
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	//!+
	// Print the results periodically.

	/**
	tick变量现在是zero value
	下面的if,如果执行了,就不是zero value

	而是否时zero value,对 select就有很重要的意义

	The zero value for a channel is nil. Perhaps surprisingly, nil channels are sometimes useful.
	Because send and receive operations on a nil channel block forever, a case in a select statement
	whose channel is nil is never selected. This lets us use nil to enable or disable cases that correspond
	to features like handling timeouts or cancellation, responding to other input events, or emitting output.
	*/
	var tick <-chan time.Time
	if *verbose {
		// 1秒(s)=1000毫秒(ms)
		tick = time.Tick(500 * time.Millisecond)
	}
	// nfiles: 总文件数, nbytes: 总字节数
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				// 注意, !ok 说明了两个事实: channel was closed and channel was drained
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			// 注意如果 tick 是 nil, 这个 case 永远不会被执行
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

//!-

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
