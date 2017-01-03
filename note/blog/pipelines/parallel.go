// +build OMIT

package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string // 文件路径
	sum  [md5.Size]byte // 文件的md5
	err  error // 计算md5过程中的错误(实际是读取文件的错误,因为md5不会返回错误)
}

// sumFiles starts goroutines to walk the directory tree at root and digest each
// regular file.  These goroutines send the results of the digests on the result
// channel and send the result of the walk on the error channel.  If done is
// closed, sumFiles abandons its work.
func sumFiles(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c.  Send the result of the walk on errc.
	c := make(chan result) // 无缓冲
	errc := make(chan error, 1) // 缓冲区长度为1, 让 send 操作不会阻塞
	go func() { // HL
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			// 启动 goroutine 之前进行 wg.Add
			wg.Add(1)
			go func() { // HL
				data, err := ioutil.ReadFile(path)
				// 要么 send result 到 c 成功
				// 要么 收到 done 完成信号
				select {
				case c <- result{path, md5.Sum(data), err}: // HL
				case <-done: // HL
				}
				// goroutine 结束之前进行 wg.Done
				wg.Done()
			}()
			
			// 每当针对一个文件启动一个 goroutine 后, 再检测是否收到 done 完成信号
			// 如果收到完成信号,返回 errors.New("walk canceled")
			// 否则,返回无错误
			// Abort the walk if done is closed.
			select {
			case <-done: // HL
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		// Walk has returned, so all calls to wg.Add are done.  Start a
		// goroutine to close c once all the sends are done.
		go func() { // HL
			// 等待所有进行 md5 的 goroutine 完成
			wg.Wait()
			// 标记不会再有 result 发送
			close(c) // HL
		}()
		// No select needed here, since errc is buffered.
		// 由于 errc 是带缓冲的 chan, 因此不会阻塞, 无需使用 select.
		errc <- err // HL
	}()
	return c, errc
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.
// inflight ['inflait] adj. 飞行中的；飞行中发生的
func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{}) // HLdone
	defer close(done)           // HLdone

	c, errc := sumFiles(done, root) // HLdone

	m := make(map[string][md5.Size]byte)
	for r := range c { // HLrange
		if r.err != nil {
			// range 从 channel 中接收到的 result 的 err 字段出错了
			return nil, r.err
		}
		// range 从 channel 中接收到的 result 的 err 字段没有出错
		m[r.path] = r.sum
	}
	// 现在,接收完毕
	
	// 检查 filepath.Walk 是否出错
	if err := <-errc; err != nil {
		// filepath.Walk 出错了
		return nil, err
	}
	return m, nil
}

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
