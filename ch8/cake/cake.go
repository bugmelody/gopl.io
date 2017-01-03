// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 234.

// Package cake provides a simulation of
// a concurrent cake shop with numerous parameters.
//
// Use this command to run the benchmarks:
// 	$ go test -bench=. gopl.io/ch8/cake
package cake

import (
	"fmt"
	"math/rand"
	"time"
)

type Shop struct {
	Verbose        bool
	Cakes          int           // number of cakes to bake
	BakeTime       time.Duration // time to bake one cake
	BakeStdDev     time.Duration // standard deviation of baking time
	// baking and icing 之间 chan 的缓冲区长度
	BakeBuf        int           // buffer slots between baking and icing
	// 多少个进行 ice 的厨师
	NumIcers       int           // number of cooks doing icing
	IceTime        time.Duration // time to ice one cake
	IceStdDev      time.Duration // standard deviation of icing time
	// icing and inscribing 之间 chan 的缓冲区长度
	IceBuf         int           // buffer slots between icing and inscribing
	InscribeTime   time.Duration // time to inscribe one cake
	InscribeStdDev time.Duration // standard deviation of inscribing time
}

// cake 代表一个蛋糕, 底层的整数是对应的标识
type cake int

func (s *Shop) baker(baked chan<- cake) {
	for i := 0; i < s.Cakes; i++ {
		c := cake(i)
		if s.Verbose {
			fmt.Println("baking", c)
		}
		work(s.BakeTime, s.BakeStdDev)
		// 蛋糕 bake 完毕,送入 baked 这个 channel
		baked <- c
	}
	// 所有蛋糕 bake 完毕, 关闭 baked 这个 channel
	close(baked)
}

func (s *Shop) icer(iced chan<- cake, baked <-chan cake) {
	// 上面 baker 函数中, 蛋糕在 bake 完了之后, 会关闭 baked channel
	// baked chan 被关闭后,下面这个 range 循环会终止
	for c := range baked {
		if s.Verbose {
			fmt.Println("icing", c)
		}
		work(s.IceTime, s.IceStdDev)
		iced <- c
	}
	// 没有必要 close iced 这个 chan, 因为在 inscriber 是通过 s.Cakes 确定的接收次数
}

func (s *Shop) inscriber(iced <-chan cake) {
	for i := 0; i < s.Cakes; i++ {
		c := <-iced
		if s.Verbose {
			fmt.Println("inscribing", c)
		}
		work(s.InscribeTime, s.InscribeStdDev)
		if s.Verbose {
			fmt.Println("finished", c)
		}
	}
}

// Work runs the simulation 'runs' times.
func (s *Shop) Work(runs int) {
	for run := 0; run < runs; run++ {
		baked := make(chan cake, s.BakeBuf)
		iced := make(chan cake, s.IceBuf)
		go s.baker(baked)
		for i := 0; i < s.NumIcers; i++ {
			go s.icer(iced, baked)
		}
		s.inscriber(iced)
	}
}

// work blocks the calling goroutine for a period of time
// that is normally distributed around d
// with a standard deviation of stddev.
func work(d, stddev time.Duration) {
	delay := d + time.Duration(rand.NormFloat64()*float64(stddev))
	time.Sleep(delay)
}
