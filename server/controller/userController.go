package controller

import (
	"encoding/json"
	"fmt"
	"net"
	"project/chatRoom/common/message"
	"project/chatRoom/common/operationData"
	"project/chatRoom/server/model"
)

type UserController struct {
	Conn net.Conn
	//用于分辨用户
	UserId int
}


//处理用户登录
func (this *UserController)ServerProcessLogin(mes *message.Message)(err error){
	//先从mes中提取data
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data),&loginMes)
	if err != nil {
		fmt.Println("login unmarshal fail err=",err)
		return
	}
	//声明一个resmes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//声明一个loginResMes
	var loginResMes message.LoginResMes
	//user Controller
	user,err:=model.MyUserModel.Login(loginMes.UserId,loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD  {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	}else{
		loginResMes.Code =200
		//到这一步用户已经登录成功，此时就把在线用户放到userManager里
		this.UserId = loginMes.UserId
		userMg.AddOnlineUser(this)
		for id,_:= range userMg.onlineUsers{
			loginResMes.OnlineUsersId = append(loginResMes.OnlineUsersId,id)
		}
		fmt.Println(user,"登陆成功")
	}


	//序列化loginResMes
	data,err :=json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("loginResMes marshal fail err=",err)
	}
	resMes.Data = string(data)
	//对resMes进行序列化
	data,err =json.Marshal(resMes)
	if err != nil {
		fmt.Println("loginResMes marshal fail err=",err)
	}
	//package操作类
	pakOper := &operationData.PackageOperation{
		Conn:this.Conn,
	}
	err = pakOper.WritePackage(data)
	if err != nil {
		fmt.Println("登录响应数据发送失败")
	}
	fmt.Println("登录响应数据发送成功")
	return
}

//处理用户注册
func (this *UserController)ServerProcessRegister(mes *message.Message)(err error){
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data),&registerMes)
	if err != nil {
		fmt.Println("registerMes 反序列化失败 err=",err)
		return
	}
	//声明一个注册响应结构体
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//处理得到的数据
	err = model.MyUserModel.Register(&registerMes.User)
	if err != nil {
		if err ==model.ERROR_USER_EXISTS {
			registerResMes.Code = 500
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		}else {
			registerResMes.Code = 600
			registerResMes.Error = "注册发生未知错误"
		}
	}else {
		registerResMes.Code = 200
	}
	data,err :=json.Marshal(registerResMes)
	//赋值给响应结构体
	resMes.Data = string(data)
	data,err = json.Marshal(resMes)
	//发送数据给客户端
	operData := operationData.PackageOperation{
		Conn:this.Conn,
	}

	err = operData.WritePackage(data)
	return
}