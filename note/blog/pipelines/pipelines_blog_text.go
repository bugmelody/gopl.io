/**
Go Concurrency Patterns: Pipelines and cancellation
[[[4-over]]] 2017-9-15 13:33:03

Introduction

Go's concurrency primitives make it easy to construct streaming data pipelines that make efficient use of I/O and multiple CPUs. This article presents examples of such pipelines, highlights subtleties that arise when operations fail, and introduces techniques for dealing with failures cleanly.

What is a pipeline?

informally [in'fɔ:məli] adv. 非正式地；不拘礼节地
There's no formal definition of a pipeline in Go; it's just one of many kinds of concurrent programs. Informally, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function. In each stage, the goroutines

    * receive values from upstream via inbound channels
    * perform some function on that data, usually producing new values
    * send values downstream via outbound channels

Each stage has any number of inbound and outbound channels, except the first and last stages, which have only outbound or inbound channels, respectively. The first stage is sometimes called the source or producer; the last stage, the sink or consumer.

sink 在计算机术语中表示 : 汇集,接收点

We'll begin with a simple example pipeline to explain the ideas and techniques. Later, we'll present a more realistic example.

Squaring numbers

Consider a pipeline with three stages.

emit [ɪ'mɪt] vt. 发出，放射；发行；发表
The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list. The gen function starts a goroutine that sends the integers on the channel and closes the channel when all the values have been sent:
 */
 
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out) // 通知下一个接收方,数据发送完毕
    }()
    return out
}

/**
The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer. After the inbound channel is closed and this stage has sent all the values downstream, it closes the outbound channel:
 */
 
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        // 既然能运行到这里,说明in已经被close
        close(out) // 通知下一个接收方,数据发送完毕
    }()
    return out
}

/**
The main function sets up the pipeline and runs the final stage: it receives values from the second stage and prints each one, until the channel is closed:
 */


func main() {
    // Set up the pipeline.
    c := gen(2, 3)
    out := sq(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}

/**
Since sq has the same type for its inbound and outbound channels, we can compose it any number of times. We can also rewrite main as a range loop, like the other stages:
 */
 
func main() {
    // Set up the pipeline and consume the output.
    for n := range sq(sq(gen(2, 3))) {
        // 2*2 = 4 , 4*4 = 16
        // 3*3 = 9 , 9*9 = 81
        fmt.Println(n) // 16 then 81
    }
}

/**
Fan-out, fan-in
Fan-out, fan-in 的意义见: http://yaotiaochimei.blog.51cto.com/4911337/861438

Multiple functions can read from the same channel until that channel is closed; this is called fan-out. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.
fan-out: 一个chan将数据发给多个worker function

A function can read from multiple inputs and proceed until all are closed by multiplexing the input channels onto a single channel that's(指那个single channel) closed when all the inputs are closed. This is called fan-in.
fan-in: 一个function从多个 input channels 读取数据.

We can change our pipeline to run two instances of sq, each reading from the same input channel. We introduce a new function, merge, to fan in the results:
 */
 
func main() {
    in := gen(2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(in)
    c2 := sq(in)

    // Consume the merged output from c1 and c2.
    for n := range merge(c1, c2) {
        fmt.Println(n) // 4 then 9, or 9 then 4
    }
}

/**
The merge function converts a list of channels to a single channel by starting a goroutine for each inbound channel that copies the values to the sole outbound channel. Once all the output(见下方源码) goroutines have been started, merge starts one more goroutine to close the outbound channel after all sends on that channel are done.

Sends on a closed channel panic, so it's important to ensure all sends are done before calling close. The sync.WaitGroup type provides a simple way to arrange this synchronization:
 */


// 返回的out在所有的 output goroutine 完成 接收,发送后(也就是cs所有的c被消耗完毕后),会被自动close掉.
// cs:待合并的channels
// 返回值,合并后的channel
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int) // 本函数最后将要返回的chan,cs参数将被合并到out中

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        // 到这里,说明 c 已经被 close 并且 消耗完毕
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        // 到这里, 所有 output goroutine 已经结束了
        close(out) // 通知下一个阶段的接受者
    }()
    return out
}
/**
Stopping short

There is a pattern to our pipeline functions:

stages close their outbound channels when all the send operations are done.
stages keep receiving values from inbound channels until those channels are closed.

This pattern(指上面的两句话) allows each receiving stage to be written as a range loop and ensures that all goroutines exit once all values have been successfully sent downstream.
[range语句会持续接收channel中的值,直到channel被close并且channel中的值被消耗完毕]

make progress: 取得进展；前进
But in real pipelines, stages don't always receive all the inbound values. Sometimes this is by design: the receiver may only need a subset of values to make progress(对于接受者来说,可能只需要一小部分值). More often, a stage exits early because an inbound value represents an error in an earlier stage. In either case the receiver should not have to wait for the remaining values to arrive, and we want earlier stages to stop producing values that later stages don't need.

indefinitely [ɪn'defɪnɪtlɪ] adv. 不确定地，无限期地；模糊地，不明确地
In our example pipeline, if a stage fails to consume all the inbound values, the goroutines attempting to send those values will block indefinitely(无限期地):
 */

    // Consume the first value from output.
    out := merge(c1, c2)
    fmt.Println(<-out) // 4 or 9
    return // 只接收了一个值就return了,另外还有一个值没有被接收
    // Since we didn't receive the second value from out,
    // one of the output goroutines is hung attempting to send it.

/**
This is a resource leak(资源泄漏): goroutines consume memory and runtime resources, and heap references in goroutine stacks keep data from being garbage collected. Goroutines are not garbage collected; they must exit on their own.

goroutine不会被垃圾回收,必须自己想办法退出.

下面开始介绍有哪些办法当 downstream stages fail to receive all the inbound values 的时候, 能让 upstream stages of our pipeline to exit.
第一种方法是使用缓冲chan.

We need to arrange for the upstream stages of our pipeline to exit even when the downstream stages fail to receive all the inbound values. One way to do this is to change the outbound channels to have a buffer. A buffer can hold a fixed number of values; send operations complete immediately if there's room in the buffer:
 */
 
c := make(chan int, 2) // buffer size 2
c <- 1  // succeeds immediately
c <- 2  // succeeds immediately
c <- 3  // blocks until another goroutine does <-c and receives 1

/**
当要发送给后续阶段的值的数量是已知的时候,用buffer chan 最合适

When the number of values to be sent is known at channel creation time, a buffer can simplify the code. For example, we can rewrite gen to copy the list of integers into a buffered channel and avoid creating a new goroutine:
 */

func gen(nums ...int) <-chan int {
    out := make(chan int, len(nums))
    for _, n := range nums {
        // 由于 out 这个 chan 的 buffer 长度为 len(nums), 因此这里的 send 操作不会阻塞
        // 由于不会阻塞,因此没有必要单独创建一个goroutine用于send
        out <- n
    }
    // 这里调用 close 会有问题吗? 参考 go doc builtin.close,
    // close 只是标记一个状态,标记此状态表明了不能再send东西到 channel 中, 因此是没有问题的
    close(out)
    return out
}


/**
Returning to the blocked goroutines in our pipeline, we might consider adding a buffer to the outbound channel returned by merge:
 */

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    /**
    // 这里为什么是1,其实是接续前面代码中,在main中只接收一次的部分,前面代码中有如下两行:
    // fmt.Println(<-out) // 4 or 9
    // return // 只接收了一个值就return了,另外还有一个值没有被接收
     */
    out := make(chan int, 1) // enough space for the unread inputs
    // ... the rest is unchanged ...

/**
While this fixes the blocked goroutine in this program, this is bad code. The choice of buffer size of 1 here depends on knowing the number of values merge will receive and the number of values downstream stages will consume. This is fragile: if we pass an additional value to gen, or if the downstream stage reads any fewer values, we will again have blocked goroutines.

Instead, we need to provide a way for downstream stages to indicate to the senders that they will stop accepting input.
相反,我们应该提供一种方式,用于后续阶段可以通知发送方:"我们停止接收了"

此时仅仅通过缓冲chan无法满足.

Explicit cancellation

When main(指从out接收值的main函数) decides to exit without receiving all the values from out, it must tell the goroutines in the upstream stages to abandon the values they're trying it send. It does so by sending values on a channel called done. It sends two values since there are potentially two blocked senders:
*/

func main() {
    in := gen(2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(in)
    c2 := sq(in)

    // Consume the first value from output.
    done := make(chan struct{}, 2)// 缓冲为2, 因为有 2 个 sq goroutine 的 worker
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // Tell the remaining senders we're leaving.
    done <- struct{}{} // 通知1个sq goroutine
    done <- struct{}{} // 通知1个sq goroutine
}

/**
The sending goroutines replace their send operation with a select statement that proceeds either when the send on out happens or when they receive a value from done(被告知无需再发送). The value type of done is the empty struct because the value doesn't matter: it is the receive event that indicates the send on out should be abandoned. The output goroutines continue looping on their inbound channel, c, so the upstream stages are not blocked. (We'll discuss in a moment how to allow this loop to return early.)
 */

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed or it receives a value
    // from done, then output calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            select {
            case out <- n:
                // 如果从c中接收到值,send到out
            case <-done:
                // 如果收到了done信号,意味着下游通知不需要发送了
                // 这里不退出,保证把源头的值拉取完
            }
        }
        // 到这里, c 已经被 close 并消耗完毕.
        wg.Done()
    }
    // ... the rest is unchanged ...
    
/**
This approach has a problem: each downstream receiver needs to know the number of potentially blocked upstream senders and arrange to signal those senders on early return. Keeping track of these counts is tedious and error-prone.

unbounded [,ʌn'baundid] adj. 1.无限的；无边的 2.不受约束的；不受控制的
We need a way to tell an unknown and unbounded number of goroutines to stop sending their values downstream. In Go, we can do this by closing a channel, because a receive operation on a closed channel can always proceed immediately, yielding the element type's zero value.

This means that main can unblock all the senders simply by closing the done channel. This close is effectively a broadcast signal to the senders. We extend each of our pipeline functions to accept done as a parameter and arrange for the close to happen via a defer statement, so that all return paths from main will signal the pipeline stages to exit.
 */
 
func main() {
    // Set up a done channel that's shared by the whole pipeline,
    // and close that channel when this pipeline exits, as a signal
    // for all the goroutines we started to exit.
    done := make(chan struct{})
    defer close(done) // 当main返回的时候通过close进行广播通知

    in := gen(done, 2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(done, in)
    c2 := sq(done, in)

    // Consume the first value from output.
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // done will be closed by the deferred call.
}


/**
Each of our pipeline stages is now free to return as soon as done is closed. The output routine in merge can return without draining its inbound channel, since it knows the upstream sender, sq, will stop attempting to send when done is closed. output ensures wg.Done is called on all return paths via a defer statement:
 */


func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c or done is closed, then calls
    // wg.Done.
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            select {
            case out <- n:
            // 如果从c中接收到值,send到out
            case <-done:
            // 说明 done 已经被 close
            // 匿名函数返回,运行defer
                return
            }
        }
    }
    // ... the rest is unchanged ...

/**
Similarly, sq can return as soon as done is closed. sq ensures its out channel is closed on all return paths via a defer statement:
 */

func sq(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-done:
                // 说明 done 已经被 close
                // 匿名函数返回,运行defer
                return
            }
        }
    }()
    return out
}

/**
Here are the guidelines for pipeline construction:

stages close their outbound channels when all the send operations are done.
stages keep receiving values from inbound channels until those channels are closed or the senders are unblocked.
Pipelines unblock senders either by ensuring there's enough buffer for all the values that are sent or by explicitly signalling senders when the receiver may abandon the channel.
 */
 
/**
Digesting a tree

Let's consider a more realistic pipeline.

MD5 is a message-digest algorithm that's useful as a file checksum. The command line utility md5sum prints digest values for a list of files.

% md5sum *.go
d47c2bbc28298ca9befdfbc5d3aa4e65  bounded.go
ee869afd31f83cbb2d10ee81b2b831dc  parallel.go
b88175e65fdcbc01ac08aaf1fd9b5e96  serial.go

Our example program is like md5sum but instead takes a single directory as an argument and prints the digest values for each regular file under that directory, sorted by path name.

% go run serial.go .
d47c2bbc28298ca9befdfbc5d3aa4e65  bounded.go
ee869afd31f83cbb2d10ee81b2b831dc  parallel.go
b88175e65fdcbc01ac08aaf1fd9b5e96  serial.go

The main function of our program invokes a helper function MD5All, which returns a map from path name to digest value, then sorts and prints the results:
 */

func main() {
    // Calculate the MD5 sum of all files under the specified directory,
    // then print the results sorted by path name.
    m, err := MD5All(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
    var paths []string
    for path := range m {
        paths = append(paths, path)
    }
    sort.Strings(paths)
    for _, path := range paths {
        fmt.Printf("%x  %s\n", m[path], path)
    }
}

The MD5All function is the focus of our discussion. In serial.go, the implementation uses no concurrency and simply reads and sums each file as it walks the tree.

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.

func MD5All(root string) (map[string][md5.Size]byte, error) {
    m := make(map[string][md5.Size]byte) // 函数最终要返回的结果
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.Mode().IsRegular() {
            return nil
        }
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }
        m[path] = md5.Sum(data) // m是外层变量
        return nil
    })
    if err != nil {
        return nil, err
    }
    return m, nil
}

/**
Parallel digestion

In parallel.go, we split MD5All into a two-stage pipeline. The first stage, sumFiles, walks the tree, digests each file in a new goroutine, and sends the results on a channel with value type result:
 */


type result struct {
    path string // 文件
    sum  [md5.Size]byte // 文件对应的md5
    err  error // ioutil.ReadFile(path) 返回的 err
}

/**
sumFiles returns two channels: one for the results and another for the error returned by filepath.Walk. The walk function starts a new goroutine to process each regular file, then checks done. If done is closed, the walk stops immediately:
 */
 
func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
    // For each regular file, start a goroutine that sums the file and sends
    // the result on c.  Send the result of the walk on errc.
    c := make(chan result) // 整个函数的第1返回值
    errc := make(chan error, 1) // 整个函数的第2返回值
    
    go func() {
        var wg sync.WaitGroup
        err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            // filepath.Walk的回调函数体内部,处理的是walk过程中的一个文件.
            // err是walk到path的错误,filepath.Walk的回调函数有权根据err决定如何处理这个错误.
            // 如果回调函数返回error,整个walk结束.
            // 如果回调函数返回err=SkipDir,表示会忽略错误.
            
            if err != nil { // 无法walk到path
                return err // 返回这个错误
            }
            if !info.Mode().IsRegular() { // 不是常规文件
                return nil
            }
            wg.Add(1)
            go func() {
                data, err := ioutil.ReadFile(path)
                // 以下select的两个case,只会有一个先发生
                select {
                case c <- result{path, md5.Sum(data), err}: // 发送sum结果到c
                case <-done: // 任务已经取消
                }
                wg.Done()
            }()
            // Abort the walk if done is closed.
            select {
            case <-done:
                return errors.New("walk canceled") // case里面可以return
            default:
                return nil // case里面可以return
            }
        })
        // Walk has returned, so all calls to wg.Add are done.  Start a
        // goroutine to close c once all the sends are done.
        go func() {
            wg.Wait()
            close(c) // 通知接收方没有数据发送了
        }()
        // No select needed here, since errc is buffered.
        errc <- err // err 是 filepath.Walk 的返回值, 将其发送到 errc
    }()
    return c, errc
}

/**
MD5All receives the digest values from c. MD5All returns early on error, closing done via a defer:
 */

func MD5All(root string) (map[string][md5.Size]byte, error) {
    // MD5All closes the done channel when it returns; it may do so before
    // receiving all the values from c and errc.
    done := make(chan struct{})
    defer close(done)

    c, errc := sumFiles(done, root)

    m := make(map[string][md5.Size]byte)
    for r := range c {
        if r.err != nil {
            return nil, r.err // 此时会运行 defer close(done)
        }
        m[r.path] = r.sum
    }
    if err := <-errc; err != nil {
        return nil, err // 此时会运行 defer close(done)
    }
    return m, nil
}

/**
Bounded parallelism(受限的并发)
bounded ['baundid] bind的过去分词 adj. 1.受限制的；有限的；狭窄的 2.【数学】有界的

The MD5All implementation in parallel.go starts a new goroutine for each file. In a directory with many large files, this may allocate more memory than is available on the machine.

We can limit these allocations by bounding the number of files read in parallel. In bounded.go, we do this by creating a fixed number of goroutines for reading files. Our pipeline now has three stages: walk the tree, read and digest the files, and collect the digests.

The first stage, walkFiles, emits the paths of regular files in the tree:
 */


func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
    paths := make(chan string) // 函数第一个返回值,所有被walk出来的path发送到这个chan
    errc := make(chan error, 1) // 函数第二个返回值,filepath.Walk返回的error; buffer选1是因为只有一次filepath.Walk的调用
    go func() {
        // Close the paths channel after Walk returns.
        defer close(paths)
        // No select needed for this send, since errc is buffered.
        errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.Mode().IsRegular() {
                return nil
            }
            select {
            case paths <- path:
            case <-done:
                return errors.New("walk canceled")
            }
            return nil
        })
    }()
    return paths, errc
}

/**
The middle stage starts a fixed number of digester goroutines that receive file names from paths and send results on channel c:
 */
 
func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
    for path := range paths {
        data, err := ioutil.ReadFile(path)
        select {
        case c <- result{path, md5.Sum(data), err}:
        case <-done:
            // 取消信号的广播
            return
        }
    }
    // 到这里,上面的for循环结束,其实说明paths已经被close并且读取完毕
}

/**
Unlike our previous examples, digester does not close its output channel, as multiple goroutines are sending on a shared channel. Instead, code in MD5All arranges for the channel to be closed when all the digesters are done:
 */
 
    // Start a fixed number of goroutines to read and digest files.
    c := make(chan result) // 所有 digester 线程公用的 chan
    var wg sync.WaitGroup
    const numDigesters = 20 // 最多同时出现20个文件读取操作
    wg.Add(numDigesters)
    for i := 0; i < numDigesters; i++ {
        go func() {
            digester(done, paths, c)
            wg.Done()
        }()
    }
    go func() {
        wg.Wait()
        // 到这里,说明所有的 digester 线程已经结束(digester函数中用了for循环)
        close(c)
    }()

/**
We could instead have each digester create and return its own output channel, but then we would need additional goroutines to fan-in the results.
我们当然也可以让每个digester自己创建并返回它的output channel,如果这么做的话,就需要一个额外的goroutine来将多个output chan进行 fan-in

The final stage receives all the results from c then checks the error from errc. This check cannot happen any earlier, since before this point, walkFiles may block sending values downstream:
 */
 
    m := make(map[string][md5.Size]byte)
    for r := range c {
        if r.err != nil {
            return nil, r.err
        }
        m[r.path] = r.sum
    }
    // Check whether the Walk failed.
    if err := <-errc; err != nil {
        // 这个检查不能在range c之前,因为如果提到range c之前,此时walkFiles可能还在阻塞发送值给downstream
        // 也就是说,造成尝试从一个buffer为1的chan中拉取值,但是chan中没有被填充的情况下,会阻塞
        return nil, err
    }
    return m, nil
}

Conclusion

This article has presented techniques for constructing streaming data pipelines in Go. Dealing with failures in such pipelines is tricky, since each stage in the pipeline may block attempting to send values downstream, and the downstream stages may no longer care about the incoming data. We showed how closing a channel can broadcast a "done" signal to all the goroutines started by a pipeline and defined guidelines for constructing pipelines correctly.

Further reading:

Go Concurrency Patterns (video) presents the basics of Go's concurrency primitives and several ways to apply them.
Advanced Go Concurrency Patterns (video) covers more complex uses of Go's primitives, especially select.
Douglas McIlroy's paper Squinting at Power Series shows how Go-like concurrency provides elegant support for complex calculations.
By Sameer Ajmani