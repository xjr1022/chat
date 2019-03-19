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