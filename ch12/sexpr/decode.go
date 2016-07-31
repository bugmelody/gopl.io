// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 344.

// Package sexpr provides a means for converting Go objects to and
// from S-expressions.
package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

//!+Unmarshal
// Unmarshal parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // get the first token

	/** Errors encountered during parsing result in a panic, so Unmarshal uses a deferred call to recover from
	the panic (§5.10) and return an error message instead. */
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

//!-Unmarshal

/**
Since a typical parser may need to inspect the current token several times, but the
Scan method advances the scanner, we wrap the scanner in a helper type called lexer
that keeps track of the token most recently returned by Scan.
*/
//!+lexer
type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

//!-lexer

/**
read function, reads the S-expression that starts with the current token and updates the
variable referred to by the addressable reflect.Value v.
*/

// The read function is a decoder for a small subset of well-formed
// S-expressions.  For brevity of our example, it takes many dubious
// shortcuts.
//
// The parser assumes
// - that the S-expression input is well-formed; it does no error checking.
// - that the S-expression input corresponds to the type of the variable.
// - that all numbers in the input are non-negative decimal integers.
// - that all keys in ((key value) ...) struct syntax are unquoted symbols.
// - that the input does not contain dotted lists such as (1 2 . 3).
// - that the input does not contain Lisp reader macros such 'x and #'x.
//
// The reflection logic assumes
// - that v is always a variable of the appropriate type for the
//   S-expression value.  For example, v must not be a boolean,
//   interface, channel, or function, and if v is an array, the input
//   must have the correct number of elements.
// - that v in the top-level call to read has the zero value of its
//   type and doesn't need clearing.
// - that if v is a numeric variable, it is a signed integer.

//!+read
func read(lex *lexer, v reflect.Value) {
	/**
	Our S-expressions use identifiers for two distinct purposes, struct field names
	and the nil value for a pointer. The read function only handles the latter case.
	When it encounters the scanner.Ident "nil", it sets v to the zero value of its
	type using the reflect.Zero function. For any other identifier, it reports an
	error. The readList function, which we’ll see in a moment, handles identifiers used as struct field names.
	*/
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if lex.text() == "nil" {
			// v.Type() : 获取 v 对应的 reflect.Type 值
			// reflect.Zero : 返回其参数对应的零值
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		/**	 A '(' token indicates the start of a list. The second function, readList, decodes a
		list into a variable of composite type — a map, struct, slice, or array — depending on what
		kind of Go variable we’re currently populating. In each case, the loop keeps parsing items
		until it encounters the matching close parenthesis, ')', as detected by the endList function.	 */
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

//!-read

//!+readlist
func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	/**	The interesting part is the recursion. The simplest case is an array. Until the closing ')' is seen, we
	use Index to obtain the variable for each array element and make a recursive call to read to populate it. As
	in many other error cases, if the input data causes the decoder to index beyond the end of the array, the
	decoder panics. A similar approach is used for slices, except we must create a new variable for each
	element, populate it, then append it to the slice.
	*/
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			/** v.Type().Elem() 说明:
			go doc reflect.Type, 它有个method : 'Elem() Type', 此方法是
			针对Array, Chan, Map, Ptr, or Slice 所用,返回对应的元素类型 */
			/**
			// New returns a Value representing a pointer to a new zero value
			// for the specified type.  That is, the returned Value's Type is PtrTo(typ).
			func New(typ Type) Value {
			*/
			/** reflect.New(xxx).Elem() 说明:
			$ go doc reflect.Value.Elem
			func (v Value) Elem() Value
			    Elem returns the value that the interface v contains or that the pointer v
			    points to. It panics if v's Kind is not Interface or Ptr. It returns the
			    zero Value if v is nil. */
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			/**
			// Append appends the values x to a slice s and returns the resulting slice.
			// As in Go, each x's value must be assignable to the slice's element type.
			func Append(s Value, x ...Value) Value {
			*/
			v.Set(reflect.Append(v, item))
		}

		/** The loops for structs and maps must parse a (key value) sublist on each iteration. For structs, the key is
		a symbol identifying the field. Analogous to the case for arrays, we obtain the existing variable for the
		struct field using FieldByName and make a recursive call to populate it.
			For maps, the key may be of any type, and analogous to the case for slices, we create a new variable,
			recursively populate it, and finally insert the new key/value pair into the map.	*/
	case reflect.Struct: // ((name value) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		/**	// MakeMap creates a new map of the specified type.
		func MakeMap(typ Type) Value { */
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			// 如果v是map,v.Type().Key() 返回该 map的key的 Type
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

//!-readlist

/**
The lexer uses the Scanner type from the text/scanner package to break an input stream into
a sequence of tokens such as comments, identifiers, string literals, and numeric literals.
The scanner’s Scan method advances the scanner and returns the kind of the next token, which
has type rune. Most tokens, like '(', consist of a single rune, but the text/scanner package
represents the kinds of the multi-character tokens Ident, String, and Int using small negative
values of type rune. Following a call to Scan that returns one of these kinds of token, the
scanner’s TokenText method returns the text of the token.
*/
