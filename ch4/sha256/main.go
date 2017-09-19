// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import "fmt"

//!+
import (
	"crypto/sha256"
	"gopl.io/ch13/equal"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	// c1 == c2 为什么是 false ??
	// 函数签名: func Sum256(data []byte) [Size]byte , 返回的是数组类型
	// If an array's element type is comparable then the array type is comparable too, so we may
	// directly compare two arrays of that type using the == operator, which reports whether all corresponding
	// elements are equal.
	//
	// byte is an alias for uint8, 因此会输出 [32]uint8
	//
	// %x是将byte的每个字节输出为两个16进制的数字
	// 2d 71 16 42 b7
	//
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}

//!-
