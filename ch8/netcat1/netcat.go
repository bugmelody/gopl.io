// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 221.
//!+

// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

/**
This program reads data from the connection and writes it to the standard out put until an
end-of-file condition or an error occurs.
 */
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

/**
注意 muxtXxx 这种模式
比如这里的 mustCopy,
还有 正则的 mustCompile
 */
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
	// 根据 io.Copy 的定义,直到遇到 EOF, io.Copy 才会返回
}

//!-
