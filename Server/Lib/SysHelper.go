package Lib

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"time"
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

func StringToByte(data string) []byte {
	return []byte(data)
}

func ByteToString(data []byte) string {
	return string(data)
}

func StringToInt(data string) (bool, string, int) {
	res, err := strconv.Atoi(data)
	if err != nil {
		return false, err.Error(), 0
	} else {
		return true, "", res
	}
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

func StringToFloat64(data string) (bool, string, float64) {
	s, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return false, err.Error(), 0
	} else {
		return true, "", s
	}
}

func Float64ToString(data float64) string {
	return strconv.FormatFloat(data, 'E', -1, 32)
}

func IntToBytes(data int) []byte {
	dataInt := int32(data)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, dataInt)
	return bytesBuffer.Bytes()
}

func BytesToInt(data []byte) int {
	bytesBuffer := bytes.NewBuffer(data)
	var dataInt int32
	binary.Read(bytesBuffer, binary.BigEndian, &dataInt)
	return int(dataInt)
}

func TimeNow() time.Time {
	return time.Now()
}

func TimeNowStr() string {
	return time.Now().Format("2006-01-02 15:04:05") // 2006-01-02 15:04:05 golang立项时间
}

func TimeStamp() int64 {
	return time.Now().Unix()
}

func TimeStampMS() int64 {
	return time.Now().UnixNano()
}

func TimeStampToStr(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}
