// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

/** apiece 英 [ə'piːs] adv. 每人；每个；各自地
There are four kinds of goroutine in this program. There is one instance apiece of
the main and broadcaster goroutines, and for each client connection there is one
handleConn and one clientWriter goroutine. */

/**
本程序中存在4中类型的goroutine
main 							1个
broadcaster 			1个

对于每个client
	handleConn			1个
	clientWriter		1个

也就是说,如果存在n个client,共计会有2n+2个goroutine(包含main goroutine)
 */

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//!+broadcaster
/**
client类型: server向每个客户端发送消息的管道, 对它能只能进行 send 操作
这里使用type定义了类型,实际是在handleConn函数的开始处被创建(针对每一个客户端).
在clientWriter中使用 range ch 不停的提取需要发送的客户端的消息并进行发送
*/
type client chan<- string // an outgoing message channel

var (
	/** 这里需要注意
	entering本身是个管道,而它的element type 是client,也是个管道. client这个管道类型代表server向客户端发送消息的
	通道,因为针对每个客户端独立创建,所以每个client代表了不同的用户.
	也就是说,entering代表client这个channel对应的用户已经连接上的消息的通道
	向entering发送消息,表示client对应的用户已经连接 */
	entering = make(chan client)
	leaving  = make(chan client)
	/** 所有客户端来的消息,通过messages管道进行广播.往messages发送消息的动作都发生在handleConn函数中.
	从messages接收消息的动作在broadcaster中的select中被处理 */
	messages = make(chan string) // all incoming client messages
)

/** The broadcaster listens on the global entering and leaving channels for announcements of
arriving and departing clients. When it receives one of these events, it updates the clients
set, and if the event was a departure, it closes the client’s outgoing message channel. The
broadcaster also listens for events on the global messages channel, to which each client sends
all its incoming messages. When the broadcaster receives one of these events, it broadcasts the
message to every connected client. */
/** broadcaster: Its local variable clients records the current set of connected clients.
The only information recorded about each client is the identity of its outgoing message channel */
func broadcaster() {
	/* 注意: 这里是使用 client 这个 channel 作为 key,代表对应的一个客户端 */
	// clients 被限制只在 broadcaster goroutine 中使用,因此对 clients 这个 map 的访问是安全的.
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// 收到从messages通道中来的数据,表示有信息需要发送给所有的client
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				/* 循环每一个client进行广播. 之后由clientWriter通过range接收发送给每个客户端 */
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

//!-broadcaster

/** The handleConn function creates a new outgoing
message channel for its client and announces the arrival of this client to the broadcaster over
the entering channel. Then it reads every line of text from the client, sending each line to the
broadcaster over the global incoming message channel, prefixing each message with the iden-
tity of its sender. Once there is nothing more to read from the client, handleConn announces
the departure of the client over the leaving channel and closes the connection. */
//!+handleConn
func handleConn(conn net.Conn) {
	/* 之前有这样一条定义: type client chan<- string // an outgoing message channel
	这里创建的ch实际上就是跟client类型对应的,只是管道方向未在此限制 */
	ch := make(chan string) // outgoing client messages

	/* handleConn creates a clientWriter goroutine for each client that receives messages broadcast
	to the client’s outgoing message channel and writes them to the client’s network connection. */
	go clientWriter(conn, ch)

	// 获取远端的地址,也就是client的地址.格式为: '127.0.0.1:64208', 'ip:端口号'
	who := conn.RemoteAddr().String()
	// 向ch这个channel发送数据,之后直接由clientWriter的goroutine进行处理,没有走广播的流程
	ch <- "You are " + who
	// 向 messages这个chan发送数据,之后由进入broadcaster的goroutine的广播流程
	messages <- who + " has arrived"
	// 将ch发送到entering这个channel,之后由broadcaster进行 clients map 的添加处理
	entering <- ch

	// bufio.NewScanner 返回的 scanner 默认是用 ScanLines 进行分隔
	// ScanLines 的源码中, 是以 `\r?\n` 这个正则作为分隔
	input := bufio.NewScanner(conn)
	for input.Scan() {
		// 从客户端收到的文字信息,需要广播给所有人,包括消息来源方的那个客户端
		// 向 messages这个chan发送数据,之后由broadcaster的goroutine进行广播
		messages <- who + ": " + input.Text()
	}

	if input.Err() != nil{
		// ...错误处理
	}
	// NOTE: ignoring potential errors from input.Err()

	// 之后broadcaster进行map的删除和channel的关闭
	leaving <- ch
	// 广播用户离开的消息,之后走广播流程
	messages <- who + " has left"
	conn.Close()
}

// clientWriter从ch中接受数据,并将收到的数据写入conn
// 注意此函数的ch参数限制了方向,只能从ch中接受数据
func clientWriter(conn net.Conn, ch <-chan string) {
	/* The client writer’s loop terminates when the broadcaster closes the channel
	after receiving a leaving notification. */
	for msg := range ch {
		// for ... range chan 循环停止的唯一条件是 chan 被 close, 并且 chan 中的已发送元素被消耗完
		// 向client写数据
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

/** The job of the main goroutine, shown below, is to listen for and accept incoming network
connections from clients. For each one, it creates a new handleConn goroutine */

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

/**
The clients map is confined to a single goroutine, the broadcaster, so it cannot be accessed concurrently.
The only variables that are shared by multiple goroutines are channels and instances of net.Conn, both
of which are concurrency safe. */
