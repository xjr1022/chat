package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"project/chat/common/message"
)

func Login(userId int, userPwd string) (err error) {
	//请求连接
	conn, err := net.Dial("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("net dial err =", err)
		return
	}
	defer conn.Close()

	//创建消息结构体
	var mes message.Message
	mes.Type = message.LoginMesType
	//创建loginMes结构体
	var loginMes = message.LoginMes{
		UserId:  userId,
		UserPwd: userPwd,
	}
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("loginMes 序列化失败")
		return
	}
	//data赋值给mes.data字段
	mes.Data = string(data)
	//将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes 序列化失败")
		return
	}

	//准备发送数据长度
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	//整形转换成切片
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	//发送长度
	_, err = conn.Write(buf[:4])
	if err != nil {
		fmt.Println("conn.write(length) fail", err)
		return
	}
	fmt.Printf("发送长度=%d 内容为=%s", len(data), string(data))

	return
}
