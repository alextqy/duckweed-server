package Lib

import (
	"net"
	"os"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

// func SetEnv(key string, Value string) error {
// 	return os.Setenv(key, Value)
// }

// func UnsetEnv(key string) error {
// 	return os.Unsetenv(key)
// }

func LocalIP() (bool, string, []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err.Error(), nil
	} else {
		var ips []string
		for _, ads := range addrs {
			if ipnet, ok := ads.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, string(ipnet.IP.String()))
				}
			}
		}
		return true, "", ips
	}
}
