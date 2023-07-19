package Lib

import (
	"fmt"
	"net"
	"time"
)

// socket服务器
func UdpServer(port string, content string, bufSize int) {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+port) // 转换地址，作为服务器使用时需要监听本机的一个端口
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr) // 启动UDP监听本机端口
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		var buf [128]byte
		len, addr, _ := conn.ReadFromUDP(buf[:])      // 读取操作会阻塞直至有数据可读取，返回值依次为读取数据长度、远端地址、错误信息
		fmt.Println(string(buf[:len]))                // 向终端打印收到的消息
		_, _ = conn.WriteToUDP([]byte(content), addr) // 写数据，返回值依次为写入数据长度、错误信息 // WriteToUDP()并非只能用于应答的，只要有个远程地址可以随时发消息
	}
}

// socket客户端
func UdpClient(ip string, port string) {
	udpAddr, err1 := net.ResolveUDPAddr("udp4", ip+":"+port) // 转换地址，作为客户端使用要向远程发送消息，这里用远程地址与端口号
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	conn, err2 := net.DialUDP("udp", nil, udpAddr) // 建立连接，第二个参数为nil时通过默认本地地址（猜测可能是第一个可用的地址，未进行测试）发送且端口号自动分配，第三个参数为远程端地址与端口号
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	go receive(conn) // 使用DialUDP建立连接后也可以监听来自远程端的数据

	for {
		_, err3 := conn.Write([]byte("naisu233~~~")) // 向远程端发送消息
		if err3 != nil {
			fmt.Println(err3.Error())
			return
		}
		time.Sleep(1 * time.Second) // 等待1秒
	}
}
func receive(conn *net.UDPConn) {
	for {
		var buf [128]byte
		len, err := conn.Read(buf[0:]) // 读取数据 // 读取操作会阻塞直至有数据可读取
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(buf[0:len]))
	}
}

// 获取广播地址
func BroadcastAddress() ([]string, string) {
	broadcastAddress := []string{}

	interfaces, err := net.Interfaces() // 获取所有网络接口
	if err != nil {
		return broadcastAddress, err.Error()
	}

	for _, face := range interfaces {
		// 选择 已启用的、能广播的、非回环 的接口
		if (face.Flags & (net.FlagUp | net.FlagBroadcast | net.FlagLoopback)) == (net.FlagBroadcast | net.FlagUp) {
			addrs, err := face.Addrs() // 获取该接口下IP地址
			if err != nil {
				return broadcastAddress, err.Error()
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok { // 转换成 IPNet { IP Mask } 形式
					if ipnet.IP.To4() != nil { // 只取IPv4的
						var fields net.IP // 用于存放广播地址字段（共4个字段）
						for i := 0; i < 4; i++ {
							fields = append(fields, (ipnet.IP.To4())[i]|(^ipnet.Mask[i])) // 计算广播地址各个字段
						}
						broadcastAddress = append(broadcastAddress, fields.String()) // 转换为字符串形式
					}
				}
			}
		}
	}

	return broadcastAddress, ""
}

// udp广播
func Broadcast(port string, content string) {
	conn, err := net.Dial("udp", "255.255.255.255:"+port)
	defer conn.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println(err.Error())
	}
}
