// 7.5.1. Caveat: An Interface Containing a Nil Pointer IsNon-Nil
// page 184

package main

import (
	"bytes"
	"fmt"
	"io"
)

/**
Consider the program below. With debug set to true, the main function collects
the output of the function f in a bytes.Buffer.
*/

// 如果是true,没有问题;如果是false,会报错
const debug = true

func main() {
	// 声明 buf 为 *bytes.Buffer, 但还没有指向具体的存储位置. 指针的 zero value 是 nil
	var buf *bytes.Buffer // ***** 解决方案是修改这行代码为: var buf io.Writer
	if debug {
		// 只有在 debug 为 true 的时候, buf 才会指向具体内存,不为 nil; 否则, defug 为 false, buf 为 nil, 不指向任何内存
		buf = new(bytes.Buffer) // enable collection of output
	}
	// buf 可能是一个空指针
	f(buf) // NOTE: subtly incorrect!
	if debug {
		fmt.Println(buf)
	}
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
}

/**
当 debug 设置为 false 的时候:
When main calls f, it assigns a nil pointer of type *bytes.Buffer to the out parameter, so the dynamic
value of out is nil. However, its dynamic type is *bytes.Buffer, meaning that out is a non-nil interface
containing a nil pointer value (Figure 7.5), so the defensive check out != nil is still true.

也就是说, out 这个 interface value,dynamic value 是 nil, dynamic type 是 *bytes.Buffer.
而 interface value 判断是否为 nil,是根据 dynamic type 是否为 nil 决定的.
由于 out 的 dynamic type 不是 nil, 因此 if out != nil 这个测试就通过了


总之:  A nil interface value, which contains no value at all, is not the same as an interface
value containing a pointer that happens to be nil.  
 */