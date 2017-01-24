// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 272.

//!+

// Package memo provides a concurrency-unsafe
// memoization of a function of type Func.
package memo

/**
A Memo instance holds the function f to memoize, of type Func, and the cache, which is a mapping from strings to results. Each result is simply the pair of results returned by a call to f—a value and an error. We’ll show several variations of Memo as the design progresses, but all will share these basic aspects.
*/

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	cache map[string]result
}

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// NOTE: not concurrency-safe!
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		// 如果 map 中不存在 key, 调用函数
		// 此时, res 是 zero value
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}

//!-


/**
由于本package的 Memo.Get 的设计不是并发安全,因此引发了memo2的修改
 */