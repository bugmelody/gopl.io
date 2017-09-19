// http://localhost:6060/doc/codewalk/sharemem/
// [[[3-over]]] 2017-9-15 10:25:39

/**
Introduction
Go's approach to concurrency differs from the traditional use of threads and shared memory.
Philosophically(adv. 哲学上；贤明地), it can be summarized:

Don't communicate by sharing memory; share memory by communicating.(不要通过共享内存来通信,应该通过通信来共享内存)

Channels allow you to pass references to data structures between goroutines.

If you consider this as passing around ownership of the data (the ability to read and write it), they
become a powerful and expressive(adj. 表现的；有表现力的) synchronization mechanism.
把 channel 看成在 goroutines 之间转移数据的拥有权(拥有权是指读写数据的能力)

In this codewalk we will look at a simple program that polls a list of URLs, checking
their HTTP response codes and periodically printing their state.
*/

/**
Poller:n. 轮询者；轮询器
poll:【计算机】顺序询问，轮询；探询
*/
package main

import (
	"log"
	"net/http"
	"time"
)

const (
	numPollers     = 2                // number of Poller goroutines to launch , 有多少个轮训者
	pollInterval   = 60 * time.Second // how often to poll each URL
	statusInterval = 10 * time.Second // how often to log status to stdout
	errTimeout     = 10 * time.Second // back-off(回退机制) timeout on error
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
	"http://www.163.com",            // 增加这个用来测试
	"http://www.baidu.com",          // 增加这个用来测试
	"http://cd.ganji.com/zhuuuuuu/", // 应该是404
}

/**
注:State type
The State type represents the state of a URL.
The Pollers send State values to the StateMonitor, which maintains a map of the
current state of each URL.
*/

// State represents the last-known state of a URL.
type State struct {
	url    string
	status string // the last-known state of a URL
}

/**
注:StateMonitor
The StateMonitor receives State values on a channel and periodically outputs
the state of all Resources being polled by the program.
*/

// StateMonitor maintains a map(指urlStatus) that stores the state of the URLs being
// polled, and prints the current state every updateInterval nanoseconds(纳秒).
// It returns a chan State to which resource state should be sent.
// 返回值类型: chan<- State, 这是一个用于发送的channel,返回的chan只能被send
func StateMonitor(updateInterval time.Duration) chan<- State {
	/**
	注:The updates channel
	The variable updates is a channel of State, on which the Poller goroutines send State values.
	This channel is returned by the function.
	这个变量是 StateMonitor 函数最终返回值
	*/
	updates := make(chan State)
	/**
	注:The urlStatus map:The variable urlStatus is a map of URLs to their most recent status.
	键是 url,值是 status
	*/
	urlStatus := make(map[string]string)
	/**
	注:The Ticker object
	A time.Ticker is an object that repeatedly sends a value on a channel at a specified interval.
	In this case, ticker triggers the printing of the current state to standard output every updateInterval nanoseconds.
	*/
	ticker := time.NewTicker(updateInterval)
	/**
	注:The StateMonitor goroutine
	StateMonitor will loop forever, selecting on two channels: ticker.C and update.
	The select statement blocks until one of its communications is ready to proceed(因为没有提供default语句).
	When StateMonitor receives a tick from ticker.C, it calls logState to print the current state.
	When it receives a State update from updates, it records the new status in the urlStatus map.
	Notice that this goroutine owns the urlStatus data structure, ensuring that it can only be accessed sequentially.
	也就是说,urlStatus这个map变量,被限制到了只能在下面的这个goroutine中被访问
	This prevents memory corruption issues that might arise from parallel reads and/or writes to a shared map.
	*/
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}

// logState prints a state map.
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}

/**
Resource type
A Resource represents the state of a URL to be polled:
the URL itself and the number of errors encountered since the last successful poll.

When the program starts, it allocates one Resource for each URL.
The main goroutine and the Poller goroutines send the Resources to each other on channels.
*/
// Resource represents an HTTP URL to be polled by this program.
type Resource struct {
	url      string
	errCount int
}

/**
The Poll method
The Poll method (of the Resource type) performs an HTTP HEAD request for the Resource's URL and
returns the HTTP response's status code. If an error occurs, Poll logs the message to standard
error and returns the error string instead.
*/

// Poll executes an HTTP HEAD request for url
// and returns the HTTP status string or an error string.
func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		// HEAD请求出错
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	// 限制,HEAD请求成功
	r.errCount = 0
	return resp.Status
}

/**
注:The Sleep method
Sleep calls time.Sleep to pause before sending the Resource to done.
The pause will either be of a fixed length (pollInterval) plus an additional
delay proportional(adj. 比例的，成比例的；相称的，均衡的) to the number of sequential errors (r.errCount).

This is an example of a typical Go idiom: a function intended to run inside a goroutine takes a channel,
upon which it sends its return value (or other indication of completed state).
*/
// Sleep sleeps for an appropriate interval (dependent on error state)
// before sending the Resource to done.
func (r *Resource) Sleep(done chan<- *Resource) {
	/* 参数: done chan<- *Resource , 这是一个用于send的 channel */
	// time.Duration(r.errCount): 这是一个类型转换
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount)) // 错误次数越多,睡得越久
	done <- r
}

/**
Poller function
Each Poller receives Resource pointers from an input channel.

In this program, the convention is that sending a Resource pointer on a channel
passes ownership of the underlying data from the sender(这里指main goroutine) to the receiver(这里指Poller goroutine).

Because of this convention, we know that no two goroutines will access this Resource at the same time.
This means we don't have to worry about locking to prevent concurrent access to these data structures.

The Poller processes the Resource by calling its Poll method.

It sends a State value to the status channel, to inform the StateMonitor of the result of the Poll.

Finally, it sends the Resource pointer to the out channel. This can be interpreted as the Poller
saying "I'm done with this Resource" and returning ownership of it to the main goroutine.

Several goroutines run Pollers, processing Resources in parallel.

in就是main中的pending
out就是main中的complete
*/
func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll() // 返回的是: status string or an error string

		status <- State{r.url, s} // sends a State value to the status channel, to inform the StateMonitor of the result of the Poll.

		// Finally, it sends the Resource pointer to the out channel. This can be interpreted as the
		// Poller saying "I'm done with this Resource" and returning ownership of it to the main goroutine.
		out <- r /* in 的参数声明是 in <-chan *Resource, 通过 range 得到的 r 类型为 *Resource */
	}
}

/**
main function
The main function starts the Poller and StateMonitor goroutines and then loops passing completed
Resources back to the pending channel after appropriate delays.
*/
func main() {
	/**
	注:
	Creating channels
	First, main makes two channels of *Resource, pending and complete.
	Inside main, a new goroutine sends one Resource per URL to pending and the
	main goroutine receives completed Resources from complete.

	The pending and complete channels are passed to each of the Poller
	goroutines, within which they are known as in and out.
	*/
	// Create our input and output channels.
	pending, complete := make(chan *Resource), make(chan *Resource)

	/**
	注:
	Initializing StateMonitor
	StateMonitor will initialize and launch a goroutine that stores the state of each Resource.
	We will look at this function in detail later.
	For now, the important thing to note is that it returns a channel of State, which is saved as status and passed to the Poller goroutines.
	*/
	// Launch the StateMonitor.
	status := StateMonitor(statusInterval)

	/**
	注:
	Launching Poller goroutines
	Now that it has the necessary channels, main launches a number of Poller goroutines, passing the channels as arguments.
	The channels provide the means of communication between the main, Poller, and StateMonitor goroutines.
	*/

	// Launch some Poller goroutines.
	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	/**
	注:
	Send Resources to pending
	To add the initial work to the system, main starts a new goroutine that allocates and sends one Resource per URL to pending.
	The new goroutine is necessary because unbuffered channel sends and receives are synchronous.
	必须要启动一个goroutine来进行send,因为pending是非缓冲chan
	That means these channel sends will block until the Pollers are ready to read from pending.
	也就是说send到pending的动作会阻塞,直到Pollers准备读取

	Were these sends performed in the main goroutine with fewer Pollers than channel sends, the program would reach a deadlock
	situation, because main would not yet be receiving from complete.
	如果相反的,send操作没有在新goroutine中,而是在main goroutine中直接进行send操作,要是碰上Pollers比较少的情况下,
	整个程序会达到死锁状态,因为没有什么东东会从complete进行接收操作,造成pending这个chan一直处于被占用状态

	Exercise for the reader: modify this part of the program to read a list of URLs from a file.
	(You may want to move this goroutine into its own named function.)
	*/

	// Send some Resources to the pending queue.
	go func() {
		for _, url := range urls {
			pending <- &Resource{url: url}
		}
	}()

	/**
	注:Main Event Loop
	When a Poller is done with a Resource, it sends it on the complete channel.
	This loop receives those Resource pointers from complete.
	For each received Resource, it starts a new goroutine calling the Resource's Sleep method.
	Using a new goroutine for each ensures that the sleeps can happen in parallel(而不会阻塞complete channel中的后来者).

	Note that any single Resource pointer may only be sent on either pending or complete at any one time.
	This ensures that a Resource is either being handled by a Poller goroutine or sleeping, but never both simultaneously.
	In this way, we share our Resource data by communicating.
	*/
	for r := range complete {
		/* 当Sleep睡完之后,会向pending这个channel发送消息,也就是说会重新进行轮训 */
		go r.Sleep(pending)
	}
	/**
	RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .
	For channels, the iteration values produced are the successive values sent on the channel until the channel is closed. If
	the channel is nil, the range expression blocks forever.
	因此,上面的 for ... range ... 是一个死循环,直到channel被close
	 */


}

/**
Conclusion
In this codewalk we have explored a simple example of using Go's concurrency primitives to share memory through commmunication.
This should provide a starting point from which to explore the ways in which goroutines and channels can be used
to write expressive and concise concurrent programs.
*/