package utils

import (
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func AppPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func AppTempFile(elem ...string) string {
	elem = append([]string{os.TempDir(), "frp-web"}, elem...)
	var tmpFile = filepath.Join(elem...)
	_ = os.MkdirAll(filepath.Dir(tmpFile), fs.ModePerm)
	return tmpFile
}

func IsIntranet(address string) (bool, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}
	ipAddr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return true, err
	}
	ip := ipAddr.IP
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true, nil
	} else if strings.HasSuffix(host, ".local") {
		return true, nil
	}
	return false, nil
}
