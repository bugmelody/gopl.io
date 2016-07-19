// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 194.

// Http3 is an e-commerce server that registers the /list and /price
// endpoints by calling (*http.ServeMux).Handle.
package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

//!+main

func main() {
	db := database{"shoes": 50, "socks": 5}
	/**
	// NewServeMux allocates and returns a new ServeMux.
	func NewServeMux() *ServeMux { ... }
	*/
	mux := http.NewServeMux()

	/**
		func (mux *ServeMux) Handle(pattern string, handler Handler)
	    Handle registers the handler for the given pattern. If a handler already
	    exists for pattern, Handle panics.
	*/

	// mux.Handle 的第二个参数要求 Handler 类型接口, 这里使用 http.HandlerFunc(db.list) 进行类型转换
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))

	/**
	// The HandlerFunc type is an adapter to allow the use of
	// ordinary functions as HTTP handlers.  If f is a function
	// with the appropriate signature, HandlerFunc(f) is a
	// Handler that calls f.
	type HandlerFunc func(ResponseWriter, *Request)

	// ServeHTTP calls f(w, r).
	func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
		f(w, r)
	}
	*/

	// mux 同时也是一个 http.Handler
	/**
		func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
	    ServeHTTP dispatches the request to the handler whose pattern most closely
	    matches the request URL.
	*/
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

//!-main

/*
//!+handlerfunc
package http

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
//!-handlerfunc
*/
