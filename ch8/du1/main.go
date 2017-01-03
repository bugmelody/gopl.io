// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 247.

//!+main

// The du1 command computes the disk usage of the files in a directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
 The ioutil.ReadDir function returns a slice of os.FileInfo — the same information that a call to os.Stat returns for a single file.
 For each subdirectory, walkDir recursively calls itself, and for each file, walkDir sends a message on the fileSizes channel.
 The message is the size of the file in bytes.
 
 The main function, shown below, uses two goroutines. The background goroutine calls walkDir for each directory specified on the
 command line and finally closes the fileSizes channel. The main goroutine computes the sum of the file sizes it receives from
 the channel and finally prints the total. 
 */

func main() {
	// Determine the initial directories.
	flag.Parse()
	
	// Args returns the non-flag command-line arguments.
	// ### func Args() []string { return CommandLine.args }
	roots := flag.Args()
	if len(roots) == 0 {
		// 如果没有从命令行提供目录,使用当前目录
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		// walk 完毕所有 root, 需要关闭 fileSizes 这个 channel. 通知接收者不再有数据发送.
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64
	// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	// /1e3 是 k, /1e6 是 m, /1e9 是 G
	// 这是大约的除法,标准是使用1024
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

//!-main

//!+walkDir

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			// 如果是目录
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			// 如果是文件,发送文件大小到 fileSizes
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

//!-walkDir

// The du1 variant uses two goroutines and
// prints the total after every file is found.
