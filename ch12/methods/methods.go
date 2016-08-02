// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 351.

// Package methods provides a function to print the methods of any value.
package methods

import (
	"fmt"
	"reflect"
	"strings"
)

//!+print
// Print prints the method set of the value x.
func Print(x interface{}) {
	v := reflect.ValueOf(x)
	t := v.Type()
	fmt.Printf("type %s\n", t)

	for i := 0; i < v.NumMethod(); i++ {
		methType := v.Method(i).Type()
		fmt.Printf("func (%s) %s%s\n", t, t.Method(i).Name,
			strings.TrimPrefix(methType.String(), "func"))
	}
}

//!-print

/**
Both reflect.Type and reflect.Value have a method called Method. Each t.Method(i)
call returns an instance of reflect.Method, a struct type that describes the name and type of
a single method. Each v.Method(i) call returns a reflect.Value representing a method
value (§6.4), that is, a method bound to its receiver. Using the reflect.Value.Call method
(which we don’t have space to show here), it’s possible to call Values of kind Func like this one,
but this program needs only its Type.
 */