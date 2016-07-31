// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 332.

/**
we use reflect.Value’s Kind method to discriminate the cases. Although there are infinitely
many types, there are only a finite number of kinds of type: the basic types Bool, String,
and all the numbers; the aggregate types Array and Struct; the reference types Chan, Func,
Ptr, Slice, and Map; Interface types; and finally Invalid, meaning no value at all.
(The zero value of a reflect.Value has kind Invalid.)
*/

/**
So far, our function treats each value as an indivisible thing with no internal structure — hence formatAtom.
For aggregate types (structs and arrays) and interfaces it prints only the type of the value,
and for reference types (channels, functions, pointers, slices, and maps), it prints the type
and the reference address in hexadecimal. This is less than ideal but still a major improvement, and
since Kind is concerned only with the underlying representation, format.Any works for named types too.
*/

// Package format provides an Any function that can format any value.
//!+
package format

import (
	"reflect"
	"strconv"
)

// Any formats any value as a string.
func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		// v.Pointer() 返回 uintptr, 代表一个指针地址
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

//!-
