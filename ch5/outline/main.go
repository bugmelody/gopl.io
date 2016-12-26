// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

/**
运行
go build gopl.io/ch5/outline
./fetch http://godoc.golangtc.com/ | ./outline
 */

// Outline prints the outline of an HTML document tree.
// outline ['aʊtlaɪn] n. 轮廓；大纲；概要；略图 vt. 概述；略述；描画…轮廓
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

//!+
func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc) // 注意,传递的是 nil slice
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack) // 输出当前 stack
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c) // 递归
	}
}

//!-

/**
Note one subtlety: although outline "pushes" an element on stack, there is no corresponding
pop. When outline calls itself recursively, the callee(被调用者) receives a copy of stack. Although the
callee(被调用者) may append elements to this slice, modifying its underlying array and perhaps even
allocating a new array, it doesn't modify the initial elements that are visible to the caller, so
when the function returns, the caller's stack is as it was before the call.
===============================
虽然outline会将元素"push"到stack末尾,却并没有对应的pop.
当outline递归调用自身的时候,callee接收到stack的一份copy.

虽然callee可能append元素到stack这个slice中, 或修改stack这个slice对应的底层数组, 或分配一个新的数组,
但callee并不会修改stack的初始元素(这些初始元素对caller永远可见) (caller拥有的stack变量永远都指向同一个底层数组,不会变动)
(因为 stack = append(stack, n.Data) 这条语句并不会影响 caller)

因此,当递归返回的时候, caller 的 stack 变量维持不变.(slice中的前几个元素没有变化,并且len没有变化).
 */