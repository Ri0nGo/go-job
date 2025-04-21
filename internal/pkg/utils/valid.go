package utils

import (
	"github.com/dlclark/regexp2"
	"net"
	"strconv"
	"strings"
)

var (
	emailCompile = regexp2.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, regexp2.None)
)

// IsValidIPv4Address 校验IPV4地址是否有效
func IsValidIPv4Address(addr string) bool {
	if len(addr) > 21 {
		return false
	}

	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return false
	}

	ipStr, portStr := parts[0], parts[1]

	// 校验 IP
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return false
	}

	// 校验端口
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return false
	}

	return true
}

func IsValidEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	ok, err := emailCompile.MatchString(email)
	if err != nil {
		return false
	}
	return ok
}
