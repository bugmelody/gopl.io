// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

/**
运行
go build gopl.io/ch5/outline
./fetch http://godoc.golangtc.com/ | ./outline
 */

// Outline prints the outline of an HTML document tree.
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
	outline(nil, doc)
}

func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}

//!-

/**
Note one subtlety: although outline "pushes" an element on stack, there is no corresponding
pop. When outline calls itself recursively, the callee receives a copy of stack. Although the
callee may append elements to this slice, modifying its underlying array and perhaps even
allocating a new array, it doesn’t modify the initial elements that are visible to the caller, so
when the function returns, the caller’s stack is as it was before the call.
 */