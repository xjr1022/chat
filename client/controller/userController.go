package controller

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"chat/client/model"
	"chat/common/message"
	"chat/common/operationData"
)

type UserController struct {

}

func (this *UserController)Login(userId int, userPwd string) (err error) {
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

		//初始化ConUser
		model.ConUser.Conn = conn
		model.ConUser.UserId = userId
		model.ConUser.UserStatus = message.UserOnline

		fmt.Println("登录成功")
		//显示在线用户
		fmt.Println("当前在线用户列表如下：")
		for _,v:=range loginResMes.OnlineUsersId {
			fmt.Println("用户id：",v)
			//完成OnlineUsers初始化
			user:=message.User{
				UserId:v,
				UserStatus:message.UserOnline,
			}
			model.OnlineUsers[v]=user
		}
		fmt.Print("\t\t")
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
func (this *UserController)Register(userId int,userPwd string,userName string)(err error){
	//请求连接
	conn, err := net.Dial("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("net dial err =", err)
		return
	}
	defer conn.Close()

	//创建消息结构体
	var mes message.Message
	mes.Type = message.RegisterMesType
	//创建loginMes结构体
	var RegisterMes  message.RegisterMes
	RegisterMes.User.UserId = userId
	RegisterMes.User.UserPwd = userPwd
	RegisterMes.User.UserName = userName
	//序列化数据
	data,err := json.Marshal(RegisterMes)
	if err != nil {
		fmt.Println("registerMes 序列化出错 err=",err)
	}
	mes.Data = string(data)
	data,err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes.data register 序列化错误 err=",err)
		return
	}
	operData :=&operationData.PackageOperation{
		Conn:conn,
	}
	err = operData.WritePackage(data)
	if err != nil {
		fmt.Println("注册信息发送错误 err=",err)
		return
	}
	mes,err = operData.ReadPackage()
	if err != nil {
		fmt.Println("ReadPack registerResMes err=",err)
		return
	}
	//反序列化
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data),&registerResMes)
	if err != nil {
		fmt.Println("registerResMes 反序列化错误 err=",err)

		return
	}
	if registerResMes.Code==200{
		fmt.Println("注册成功")
		os.Exit(0)
	}else{
		fmt.Println("注册失败，err=",registerResMes.Error)
		os.Exit(0)
	}
	return
}

