// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 136.

// The toposort program prints the nodes of a DAG in topological order.
package main

import (
	"fmt"
	"sort"
)

/**
The prerequisites are given in the prereqs table below, which is a mapping from
each course to the list of courses that must be completed before it.
课程 : 在学习这门课程之前需要完成的其他课程
 */

/**
This kind of problem is known as topological sorting. Conceptually, the prerequisite
information forms a directed graph with a node for each course and edges from each course to
the courses that it depends on. The graph is acyclic: there is no path from a course that leads
back to itself. We can compute a valid sequence using depth-first search through the graph
with the code below
 */

//!+table
// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

//!+main
func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	// 最后要返回的结果
	var order []string
	// seen 表示是否已经被 visitAll 访问过
	seen := make(map[string]bool)
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				// 递归调用,最深的递归最先进行 order = append(order, item) 操作, 也就是最先将最基础的学科进行 append
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}

//!-main
