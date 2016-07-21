// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 223.

// Netcat is a simple read/write client for TCP servers.
package main


/**
While the main goroutine reads the standard input and sends it to the server, a second
goroutine reads and prints the server’s response. When the main goroutine encounters the
end of the input, for example, after the user types Control-D (^D) at the terminal (or the
equivalent Control-Z on Microsoft Windows), the program stops, even if the other goroutine
still has work to do. (We’ll see how to make the program wait for both sides to finish once
we’ve introduced channels in Section 8.4.1.)
 */
import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 另起一个goroutine,用于接收server的信息到stdout
	// 如果 main goroutine 中的 mustCopy 遇到 EOF 返回, 另一个 goroutine 中仍然在运行, 造成程序退出而第二个 goroutine 被强制中断
	go mustCopy(os.Stdout, conn)
	
	// 接收stdin的信息发送到server
	mustCopy(conn, os.Stdin)
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
