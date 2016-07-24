// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 251.

// The du4 command computes the disk usage of the files in a directory.
package main

// The du4 variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/** we create a cancellation channel(也就是 done) on which no values are ever sent, but whose
closure indicates that it is time for the program to stop what it is doing. We also define a
utility function, cancelled, that checks or polls the cancellation state at the instant it is called.
*/
// 如果chan的元素类型为struct{},表示传什么值不重要,重要的是传递的这个事件
//!+1
// 注意,我们不会往done里面放任何数据,它只是用于cancellation
var done = make(chan struct{})

func cancelled() bool {
	// 对于 select,如果所有case都被阻塞,这时才会执行default
	select {
	case <-done:
		// receive成功,其实是代表done已经被close
		return true
	default:
		// 说明done没有被close
		return false
	}
}

//!-1

func main() {
	// Determine the initial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	//!+2
	/** create a goroutine that will read from the standard input, which is typically connected
	to the terminal. As soon as any input is read (for instance, the user presses the return
	key), this goroutine broadcasts the cancellation by closing the done channel. */

	// Cancel traversal when input is detected.
	// 启动一个后台goroutine,一旦侦测到用户在stdin中输入任何东西,则取消整个程序的运行(通过close(done))
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		/** Recall that after a channel has been closed and drained of all sent values, subsequent receive
		operations proceed immediately, yielding zero values. We can exploit this to create a broadcast
		mechanism: don’t send a value on the channel, close it. */
		// 发送cancel事件
		close(done)
	}()
	//!-2

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}

	// 后台单独启动一个goroutine等待WaitGroup完成,之后关闭fileSizes
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	//!+3
	for {
		select {
		case <-done:
			/** returns if this case is ever selected, but before it returns it must first drain the fileSizes
			channel, discarding all values until the channel is closed. It does this to ensure that any active
			calls to walkDir can run to completion without getting stuck sending to fileSizes */
			// Drain fileSizes to allow existing goroutines to finish.
			// 注意:fileSizes这个chan是无缓冲的,因此必须先消耗光里面的内容;否则,其他准备send的goroutine会阻塞
			for range fileSizes {
				// for range channel 循环完成的条件是channel被关闭并拉取完所有已发送的值
				// Do nothing.
			}
			// 现在 fileSizes 已经被关闭并拉取完所有已发送的值
			return
		case size, ok := <-fileSizes:
			// ...
			//!-3
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

/** The walkDir goroutine polls the cancellation status when it begins, and returns without doing
anything if the status is set. This turns all goroutines created after cancellation into no-ops: */
// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+4
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		// 关键部位判断cancelled已经很有用了
		return
	}
	for _, entry := range dirents(dir) {
		/**
				It might be profitable to poll the cancellation status again within walkDir’s loop, to avoid cre-
		ating goroutines after the cancellation event. Cancellation involves a trade-off; a quicker
		response often requires more intrusive changes to program logic. Ensuring that no expensive
		operations ever occur after the cancellation event may require updating many places in your
		code, but often most of the benefit can be obtained by checking for cancellation in a few
		important places.
		我们也可以在这个循环内加入是否cancelled的判断进行更加快速的返回避免无用goroutine的创建(go walkDir(xxx)的递归).
		但是cancellation通常都是存在trade-off,更快速的响应(更好的性能)需要修改更多的代码,让代码变得复杂.
		一般来说,在关键部分加入cancellation通常已经很有效了,没必要过分追求完美 */
		// ...
		//!-4
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			// 递归的goroutine创建
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
		//!+4
	}
}

//!-4

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
//!+5
func dirents(dir string) []os.FileInfo {
	/** A little profiling of this program revealed that the bottleneck was the acquisition of a sema-
	phore token in dirents. The select below makes this operation cancellable and reduces the
	typical cancellation latency of the program from hundreds of milliseconds to tens
	也就是说,降低了cancellation的延迟 */
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	// ...read directory...
	//!-5

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	// 注意: du3 中的 ioutil.ReadDir 内部实际是调用了 func (f *File) Readdir(n int) (fi []FileInfo, err error) {
	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
