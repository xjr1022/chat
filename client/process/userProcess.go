package process

import (
	"encoding/json"
	"fmt"
	"net"
	"project/chatRoom/common/message"
	"project/chatRoom/common/operationData"
)

type UserProcess struct {

}

func (this *UserProcess)Login(userId int, userPwd string) (err error) {
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

	packOper:=&operationData.PackageOperation{
		Conn:conn,
	}
	//发送数据
	err = packOper.WritePackage(data)

	if err != nil {
		fmt.Println("conn.write(data) fail", err)
		return
	}
	fmt.Printf(" 内容为=%s \n", string(data))

	mes,err = packOper.ReadPackage()
	fmt.Printf(" 读取登录响应数据=%s \n", mes)
	if err != nil {
		fmt.Println("loginResMes解读错误，err=",err)
		return
	}
	//将消息反序列化
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data),&loginResMes)
	if loginResMes.Code==200 {
		fmt.Println("登录成功")

		//开启一个进程保持和服务器的通讯
		go ProcessServerMes(conn)

		for{
			ShowMenu()
		}
	}else{
		fmt.Println("登陆失败，err=",loginResMes.Error)
	}

	return
}

