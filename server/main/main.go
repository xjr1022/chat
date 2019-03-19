package main

import (
	"fmt"
	"net"
	"project/chatRoom/server/model"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()
	processor :=Processor{
		Conn:conn,
	}
	processor.MainProcess(conn)

}
//完成对userModel的初始化
func initUserModel()  {

	model.MyUserModel = &model.UserModel{
		Pool:pool,
	}
}
func main() {

	//初始化redis
	initPool("localhost:6379",16,0,100*time.Second)
	initUserModel()
	//提示信息
	fmt.Println("服务器正在监听")

	listen,err := net.Listen("tcp", "0.0.0.0:1234")

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
