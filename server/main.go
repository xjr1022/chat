package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"project/chat/common/message"
)

func process(conn net.Conn) {
	defer conn.Close()

	//循环处理客户端发送的消息
	for {
		mes, err := readPackage(conn)

		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端退出")
			}else{
				fmt.Println("conn.read err=", err)

			}
			return
		}
		fmt.Println("读取到的buf=", buf)
	}
}
func main() {
	//提示信息
	fmt.Println("服务器正在监听")
	listen, err := net.Listen("tcp", "0.0.0.0:1234")
	defer listen.Close()
	if err != nil {
		fmt.Println("服务器端口监听异常 err=", err)
		return
	}
	//一旦监听成功，就等待客户端来连接
	for {
		fmt.Println("等待客户端连接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen监听成功，但发生错误 err=", err)
		}
		//连接成功，开启一个协程和客户端保持通讯
		go process(conn)
	}

}
func readPackage(conn net.Conn) (mes message.Message,err error){
	buf := make([]byte, 8096)
	fmt.Printf("读取客户端发送的数据")
	//conn.read 在conn没有被关闭的情况下才会被阻塞，先读取字节长度
	_,err = conn.Read(buf[:4])
	if err != nil{
		return
	}
	var packageLen uint32
	//转成uint32类型
	packageLen = binary.BigEndian.Uint32(buf[:4])
	n,err := conn.Read(buf[:packageLen])
	if n!= int(packageLen) || err != nil{
		return
	}
	err = json.Unmarshal(buf[:packageLen],&mes)
	if err != nil {
		fmt.Printf("json.unmarsha err=%v",err )
		return
	}
	return
}
