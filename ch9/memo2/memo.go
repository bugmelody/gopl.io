// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 275.

// Package memo provides a concurrency-safe memoization a function of
// type Func.  Concurrent requests are serialized by a Mutex.
package memo

/**
 The simplest way to make the cache concurrency-safe is to use monitor-based
 synchronization. All we need to do is add a mutex to the Memo, acquire the
 mutex lock at the start of Get, and release it before Get returns, so that
 the two cache operations occur within the critical section: 
 */

import "sync"

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

//!+

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]result
}

// Get is concurrency-safe.
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

//!-


/**
 Now the race detector is silent, even when running the tests concurrently. Unfortunately
 this change to Memo reverses our earlier performance gains. By holding the lock for the
 duration of each call to f, Get serializes all the I/O operations we intended to parallelize.
 What we need is a non-blocking cache, one that does not serialize calls to the function it memoizes.
 由此引发了 memo3 的修改.
 */