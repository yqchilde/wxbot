package net

import (
	"errors"
	"net"
	"regexp"
	"time"
)

var clientIP = "192.168.0.1"

// GetIPNet 获取IPNet对象
func GetIPNet() (*net.IPNet, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iFace := range iFaces {
		addrs, err := iFace.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					return ip, nil
				}
			}
		}
	}
	return nil, errors.New("failed get ip net")
}

// GetIPWithLocal 获取本地内网IP
func GetIPWithLocal() (string, error) {
	ipNet, err := GetIPNet()
	if err != nil {
		return clientIP, err
	}

	return ipNet.IP.String(), nil
}

// CheckoutIpPort 检出IP和端口
func CheckoutIpPort(str string) (string, error) {
	reg, err := regexp.Compile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}:[0-9]{1,5}")
	if err != nil {
		return "", errors.New("匹配Ip+Port失败")
	}
	return reg.FindString(str), nil
}

func PingConn(ipPort string, timeout time.Duration) bool {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	conn, err := net.DialTimeout("tcp", ipPort, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
