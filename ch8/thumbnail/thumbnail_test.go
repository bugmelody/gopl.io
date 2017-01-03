// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// This file is just a place to put example code from the book.
// It does not actually run any code in gopl.io/ch8/thumbnail.

package thumbnail_test

import (
	"log"
	"os"
	"sync"

	"gopl.io/ch8/thumbnail"
)

//!+1
// makeThumbnails makes thumbnails of the specified files.
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

//!-1

//!+2
// NOTE: incorrect!
//这是错误的设计,因为用go了之后,直接返回了,最后可能还在执行图像处理,主线程却已经退出了
/**
This version runs really fast—too fast, in fact, since it takes less time than the original, even
when the slice of file names contains only a single element. If there's no parallelism, how can
the concurrent version possibly run faster? The answer is that makeThumbnails returns before it
has finished doing what it was supposed to do. It starts all the goroutines, one per file name, but
doesn't wait for them to finish.
*/
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f) // NOTE: ignoring errors
	}
}

//!-2

//!+3
// makeThumbnails3 makes thumbnails of the specified files in parallel.
func makeThumbnails3(filenames []string) {
	// ch 作用是当一个 goroutine 完成后发送事件通知
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			// struct{}是类型, struct{}{}是literal
			ch <- struct{}{}
		}(f)
	}

	// 实际上是等待len(filenames)次接收
	// Wait for goroutines to complete.
	for range filenames {
		<-ch
	}
}

//!-3

//!+4
// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error if any step failed.
// 错误的设计,如果中间发生err,main goroutine return 退出,剩余未完成的 goroutine 由于 send 卡住造成无法退出.
func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil
}

//!-4

//!+5
// makeThumbnails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order,
// or an error if any step failed.
func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	// 通常定义一个struct来保存单个输入与对应的err
	type item struct {
		thumbfile string
		err       error
	}

	// 通过 buffer chan 来确保 send 不会被卡住
	ch := make(chan item, len(filenames))
	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			// 注意,这里使用了buffered chan,所以提前返回不会有问题,由于有buffer,不会造成发送方阻塞,当前还在运行中的goroutine会在send之后正常退出
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

//!-5

//!+6
// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.
// filenames参数: 接收的channel,每个元素代表一个文件名
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for f := range filenames {
		// increments WaitGroup counter
		/** Add, which increments the counter, must be called before the worker goroutine
		starts, not within it; otherwise we would not be sure that the Add happens before
		the "closer" goroutine calls Wait.
		*/
		wg.Add(1)
		// worker
		go func(f string) {
			// decrements the WaitGroup counter,放在defer语句中,确保一定被执行
			// We use defer to ensure that the counter is decremented even in the error case
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	go func() {
		wg.Wait()
		// 如果不close,主线程的的range(在下方)就根本停不下来
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}

	/**
	Observe how we create a closer goroutine that
	waits for the workers to finish before closing the sizes channel. These two operations, wait
	and close, must be concurrent with the loop over sizes. Consider the alternatives: if the wait
	operation were placed in the main goroutine before the loop, it would never end(chan中有数据却没人接收), and if placed
	after the loop, it would be unreachable since with nothing closing the channel, the loop would
	never terminate(range loop 要停止的条件是chan被close并且已发送的数据被消费完).

	*/

	return total
}

// The structure of the code above is a common and idiomatic pattern for looping
// in parallel when we don’t know the number of iterations.

//!-6
