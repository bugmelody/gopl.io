package snippet

import (
	"net/http"
)

// 已经有了一个URL,如何修改此URL的querystring
func updateQueryString(req *http.Request) {
	q := req.URL.Query()
	q.Set("q", "要设置的q参数")
	req.URL.RawQuery = q.Encode()
}
