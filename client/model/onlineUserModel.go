package model

import (
	"fmt"
	"net"
	"project/chatRoom/common/message"
)

type ConnectUser struct {
	Conn net.Conn
	message.User
}

//客户端维护的用户map
var OnlineUsers map[int]message.User = make(map[int]message.User)
//登录成功后完成对ConnectUser的初始化
var ConUser ConnectUser

//显示在线用户
func OutPutOnlineUser(){
	//遍历一把 onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range OnlineUsers{
		//如果不显示自己.
		fmt.Println("用户id:\t", id)
	}
}

//更新用户的在线列表
func UpdateUserStatus(mes message.NotifyUserStatusMes)  {

	user,ok := OnlineUsers[mes.UserId]
	if !ok {
		user = message.User{
			UserId:mes.UserId,
		}
	}
	user.UserStatus = mes.Status
	OnlineUsers[mes.UserId] = user

	OutPutOnlineUser()
}