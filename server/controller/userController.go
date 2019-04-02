package controller

import (
	"encoding/json"
	"fmt"
	"net"
	"chat/common/message"
	"chat/common/operationData"
	"chat/server/model"
)

type UserController struct {
	Conn net.Conn
	//用于分辨用户
	UserId int
}

//通知客户端用户上线
func (this *UserController) NotifyOthersUserOnline(userId int){
	//遍历OnlineUsers，然后发送
	for id,userCont :=range userMg.onlineUsers{
		if id == userId {
			continue
		}
		//同时用户上线
		userCont.NotifyUserOnline(userId)
	}
}

//根据userId发送上线通知
func (this *UserController) NotifyUserOnline(userId int){
	//实例化消息结构体
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	//实例化OnlineMes
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//序列化notifyUser
	data,err:=json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Print("notifyUser 序列化失败 err=",err)
	}
	mes.Data = string(data)

	data,err = json.Marshal(mes)
	if err != nil {
		fmt.Print("notify Mes 序列化失败 err=",err )
		return
	}

	packOper := &operationData.PackageOperation{
		Conn:this.Conn,
	}
	//发送消息给客户端
	err = packOper.WritePackage(data)
	if err != nil {
		fmt.Print("发送在线消息失败 err=",err)
		return
	}
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
		this.NotifyOthersUserOnline(loginMes.UserId)
		//将当前在线用户的id放入响应消息里
		for id,_:= range userMg.onlineUsers{
			//去掉自身的id
			if this.UserId==id {
				continue
			}
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