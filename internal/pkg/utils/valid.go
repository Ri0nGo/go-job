package utils

import (
	"net"
	"strconv"
	"strings"
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
