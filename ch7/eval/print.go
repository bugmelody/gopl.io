// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package eval

import (
	"bytes"
	"fmt"
)

// Format formats an expression as a string.
// It does not attempt to remove unnecessary parens.
func Format(e Expr) string {
	var buf bytes.Buffer
	// write 函数在下方定义
	write(&buf, e)
	/**
	String returns the contents of the unread portion of the buffer
	as a string.  If the Buffer is a nil pointer, it returns "<nil>".
	func (b *Buffer) String() string {
	*/
	return buf.String()
}

func write(buf *bytes.Buffer, e Expr) {
	switch e := e.(type) {
	case literal:
		fmt.Fprintf(buf, "%g", e)

	case Var:
		fmt.Fprintf(buf, "%s", e)

	case unary:
		fmt.Fprintf(buf, "(%c", e.op)
		write(buf, e.x)
		/**
		WriteByte appends the byte c to the buffer, growing the buffer as needed.
		The returned error is always nil, but is included to match bufio.Writer's
		WriteByte. If the buffer becomes too large, WriteByte will panic with
		ErrTooLarge.
		func (b *Buffer) WriteByte(c byte) error {
		*/
		buf.WriteByte(')')

	case binary:
		buf.WriteByte('(')
		write(buf, e.x)
		fmt.Fprintf(buf, " %c ", e.op)
		write(buf, e.y)
		buf.WriteByte(')')

	case call:
		fmt.Fprintf(buf, "%s(", e.fn)
		for i, arg := range e.args {
			if i > 0 {
				buf.WriteString(", ")
			}
			write(buf, arg)
		}
		buf.WriteByte(')')

	default:
		panic(fmt.Sprintf("unknown Expr: %T", e))
	}
}
