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
		// 在这里,对于调用者来说,stack这个sliceHeader结构体是不会变的,甚至通过stack可窥见的数组元素也不会变(可能stack指向的底层数据已经被增删)
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

stack
1,2,3
开始递归调用stack
1,2,3,4
1,2,3,4,5
1,2,3,4,5,6
1,2,3,4,5
1,2,3,4
1,2,3 回到调用者

        
 */