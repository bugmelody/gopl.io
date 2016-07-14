// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 99.

// Graph shows how to use a map of maps to represent a directed graph.
package main

import "fmt"

/**
摘自书上说明:
The addEdge function shows the idiomatic way to populate a map lazily, that is, to initialize
each value as its key appears for the first time. The hasEdge function shows how the zero
value of a missing map entry is often put to work: even if neither from nor to is present,
graph[from][to] will always give a meaningful result.
 */

//!+
var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		// if edges == nil { ::: map 对应的 zero value 是 nil, 这里实际是判断 graph[from] 对应的 map 值是否已经被初始化
		edges = make(map[string]bool)
		graph[from] = edges
	}
	
	// 现在, graph[from] 肯定已经存在, 也即是 edges 
	edges[to] = true
}

func hasEdge(from, to string) bool {
	// 如果 from 不存在, graph[from] 返回 nil (map对应的zero value)
	// 而在 nil 上面索引 to, nil[to], 是一个合法的操作
	return graph[from][to]
}

//!-

func main() {
	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))

}
