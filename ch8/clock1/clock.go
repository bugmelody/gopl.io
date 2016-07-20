// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 219.
//!+

/**
$ go build gopl.io/ch8/clock1
$ ./clock1 &

$ nc localhost 8000
21:09:56
21:09:57
^C

$ telnet localhost 8000
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
21:10:46
21:10:47
21:10:48
^C

 */

// Clock1 is a TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		/**
		// Accept waits for and returns the next connection to the listener.
		Accept() (Conn, error)
		*/
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	/**
	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error
	*/
	defer c.Close()

	for {
		/**
		// WriteString writes the contents of the string s to w, which accepts a slice of bytes.
		// If w implements a WriteString method, it is invoked directly.
		func WriteString(w Writer, s string) (n int, err error)
		*/
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

//!-
