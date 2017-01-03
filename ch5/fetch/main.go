// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 148.

// Fetch saves the contents of a URL into a local file.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

/**
It's tempting to use a second deferred call, to f.Close, to close the local file, but this would be subtly wrong because
os.Create opens a file for writing, creating it as needed(os.Create是按需创建文件). On many file systems, notably(adv. 显著地；尤其) NFS,
write errors are not reported immediately but may be postponed until the file is closed. Failure
to check the result of the close operation could cause serious data loss to go unnoticed(go unnoticed:不被差觉的).
However, if both io.Copy and f.Close fail, we should prefer to report the error from
io.Copy since it occurred first and is more likely to tell us the root cause.
 */

//!+
// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	// path.Base 文档中提到: If the path consists entirely of slashes, Base returns "/".
	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	// resp.Body 类型为 io.ReadCloser
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	if closeErr := f.Close(); err == nil {
		// 如果 io.Copy 返回的 err 不是 nil, 则优先返回 io.Copy 的错误
		// 否则,如果 io.Copy 没有出错, 则将 err 变量设置为 f.Close() 的返回值
		err = closeErr
	}

	return local, n, err
}

//!-

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}
