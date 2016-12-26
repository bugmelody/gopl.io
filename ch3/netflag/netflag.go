// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 77.

// Netflag demonstrates an integer type used as a bit field.
package main

import (
	"fmt"
	// 使用 . "net", 之后可以将 net.XXX 直接写为 XXX
	. "net"
)

//!+
// 测试 FlagUp 这个 bit 是否为 1
func IsUp(v Flags) bool     { return v&FlagUp == FlagUp }
// 将 FlagUp 这个 bit 清0
func TurnDown(v *Flags)     { *v &^= FlagUp }
// 将 FlagBroadcast 这个 bit 设置为 1
func SetBroadcast(v *Flags) { *v |= FlagBroadcast }
// 测试 FlagBroadcast, FlagMulticast 其中之一是否被设置
func IsCast(v Flags) bool   { return v&(FlagBroadcast|FlagMulticast) != 0 }

func main() {
	var v Flags = FlagMulticast | FlagUp
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10001 true"
	TurnDown(&v)
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10000 false"
	SetBroadcast(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10010 false"
	fmt.Printf("%b %t\n", v, IsCast(v)) // "10010 true"
}

//!-
