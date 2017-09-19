package snippet

import (
	"fmt"
	"net"
	"net/http"
)

// 从http请求中获取用户ip
func GetUserIpFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		// SplitHostPort error
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		// net.ParseIP error
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}
