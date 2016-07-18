// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 160.

/**
运行
$ go run ch6/urlvalues/main.go 1> /dev/null
$ go run ch6/urlvalues/main.go 2> /dev/null
 */
// The urlvalues command demonstrates a map type with methods.
package main

/*
//!+values
package url

// Values maps a string key to a list of values.
type Values map[string][]string

// Get returns the first value associated with the given key,
// or "" if there are none.
func (v Values) Get(key string) string {
	if vs := v[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

// Add adds the value to key.
// It appends to any existing values associated with key.
func (v Values) Add(key, value string) {
	// 注意这里, v[key] 可能会返回nil,但 v[key] = append(v[key], value) 仍然合法
	v[key] = append(v[key], value)
}
//!-values
*/

import (
	"fmt"
	"net/url"
)

func main() {
	//!+main
	m := url.Values{"lang": {"en"}} // direct construction
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Println(m.Get("lang")) // "en"
	fmt.Println(m.Get("q"))    // ""
	fmt.Println(m.Get("item")) // "1"      (first value)
	fmt.Println(m["item"])     // "[1 2]"  (direct map access)

	m = nil
	fmt.Println(m.Get("item")) // ""
	m.Add("item", "3")         // panic: assignment to entry in nil map
	//!-main
}

/**
In the final call to Get, the nil receiver behaves like an empty map. We could equivalently have
written it as Values(nil).Get("item")), but nil.Get("item") will not compile because
the type of nil has not been determined. By contrast, the final call to Add panics as it tries to
update a nil map.
一个nil map是不指向任何内容的,也就是没有指向任何hash table.

Because url.Values is a map type and a map refers to its key/value pairs indirectly, any
updates and deletions that url.Values.Add makes to the map elements are visible to the
caller. However, as with ordinary functions, any changes a method makes to the reference
itself, like setting it to nil or making it refer to a different map data structure, will not be
reflected in the caller.
 */