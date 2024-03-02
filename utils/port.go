package utils

import (
	"fmt"
	"net"
)

func IWantUseHttpPort(port ...int) int {
	if len(port) > 1 {
		panic("端口参数错误")
	}
	if len(port) == 0 {
		port = []int{8000}
	}
	if len(port) == 0 || port[0] <= 1024 || port[0] > 60000 {
		port = []int{8000}
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port[0]))
	if err != nil {
		return IWantUseHttpPort(port[0] + 1)
	}
	defer func() { _ = listener.Close() }()
	return port[0]
}
