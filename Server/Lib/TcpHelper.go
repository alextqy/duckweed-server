package lib

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"time"
)

// tcp socket 服务器端
func TcpServer(ip string, port string, content string) {
	addr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err.Error())
	}
	tcpListen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("tcp server on:", addr.String())
	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			fmt.Println("abnormal network monitoring", err.Error())
			continue
		}
		clientMessage(conn)
		conn.Close()
		fmt.Println(conn.RemoteAddr(), ": actively disconnect")
	}
}
func clientMessage(conn net.Conn) {
	buf := [...]byte{}
	for {
		readSize, err1 := conn.Read(buf[0:])
		if err1 != nil {
			log.Fatal(err1.Error())
		}
		remoteAddr := conn.RemoteAddr()
		fmt.Println("client:", remoteAddr, " message:", string(buf[0:]))
		_, err2 := conn.Write([]byte(string(buf[0:readSize]) + " " + time.Now().String()))
		if err2 != nil {
			log.Fatal(err2.Error())
		}
	}
}

func TcpServerPlus(ip string, port string) {
	addr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err.Error())
	}
	tcpListen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("tcp server on:", addr.String())
	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			continue
		}
		// 编程高性能并发服务器端
		go handlerConn(conn)
	}
}

func handlerConn(conn net.Conn) {
	defer fmt.Println("The client " + conn.RemoteAddr().String() + " actively disconnects")
	defer conn.Close() // 正常链接情况下，handlerConn不会释放出来到这里 当客户端强制断开，才会return到这里关闭当前conn
	fmt.Println("New client " + conn.RemoteAddr().String())

	// 获取客户端信息info 返回info服务器时间

	var buf [1024]byte
	for {
		readSize, err := conn.Read(buf[0:])
		if err != nil {
			log.Fatal(err.Error())
		}
		remoteAddr := conn.RemoteAddr()
		gorid := GetGoroutineID()
		fmt.Println("协程id: "+strconv.FormatUint(gorid, 10)+" 来自远程ip:", remoteAddr, " 的消息:", string(buf[0:readSize]))
		_, err2 := conn.Write([]byte(string(buf[0:readSize]) + " " + time.Now().String()))
		// 一定要执行下面的return 才能监听到客户端主动断开，服务器端对本次conn进行close处理 dealErrorWithReturn不能达到这个效果。
		if err2 != nil {
			return
		}
	}
}
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
