// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 278.

// Package memo provides a concurrency-safe non-blocking memoization
// of a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a monitor goroutine.
package memo

//!+Func

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

/**
the Memo type now consists of a channel, requests, through which the caller
of Get communicates with the monitor goroutine. The element type of the channel
is a request. Using this structure, the caller of Get sends the monitor goroutine
both the key, that is, the argument to the memoized function, and another
channel, response, over which the result should be sent back when it becomes
available. This channel will carry only a single value.
*/
type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

/**
The Get method, creates a response channel, puts it in the request, sends it to the monitor goroutine, then immediately receives from it.
??哪个是monitor goroutine??
*/
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

// server方法就是monitor goroutine
//!+monitor

func (memo *Memo) server(f Func) {
	// The cache variable is confined to the monitor goroutine (*Memo).server
	cache := make(map[string]*entry)

	/**
	The monitor reads requests in a loop until the request channel is closed by the Close method. For each
	request, it consults the cache, creating and inserting a new entry if none was found.
	*/
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			/**
			In a similar manner to the mutex-based version, the first request for a given key becomes
			responsible for calling the function f on that key, storing the result in the entry, and
			broadcasting the readiness of the entry by closing the ready channel. This is done by (*entry).call.
			*/
			go e.call(f, req.key) // call f(key)
		}
		/**
		A subsequent request for the same key finds the existing entry in the map, waits for the
		result to become ready, and sends the result through the response channel to the client
		goroutine that called Get. This is done by (*entry).deliver. The call and deliver methods
		must be called in their own goroutines to ensure that the monitor goroutine does not stop
		processing new requests.
		*/
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition.
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

//!-monitor

/**
This example shows that it’s possible to build many concurrent structures using either
of the two approaches—shared variables and locks, or communicating sequential processes—without
excessive complexity.

It’s not always obvious which approach is preferable in a given situation, but it’s
worth knowing how they correspond. Sometimes switching from one approach to the other
can make your code simpler.
*/
