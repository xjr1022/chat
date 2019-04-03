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

//显示菜单
func ShowMenu() {

	fmt.Println("-------恭喜xxx登录成功---------")
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 信息列表---------")
	fmt.Println("-------4. 退出系统---------")
	fmt.Println("请选择(1-4):")

	var key int
	var content string

	//聊天会经常用到ChatMesController，就定义在外面
	chaMesCont := ChatMesController{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		model.OutPutOnlineUser()
	case 2:
		fmt.Scanf("%s\n", &content)
		chaMesCont.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出了系统...")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确..")
	}
}

//用线程开启此函数，和服务器后台保持通讯
func ProcessServerMes(conn net.Conn) {
	//创建一个OperationData实例
	operDa := &operationData.PackageOperation{
		Conn: conn,
	}
	for {
		fmt.Println("等待读取客户端发送的消息")
		mes, err := operDa.ReadPackage()
		if err != nil {
			fmt.Println("协程收取数据错误 err=", err)
			return
		}
		//处理服务器发来的消息
		switch mes.Type {
		case message.NotifyUserStatusMesType: //用户上线
			var notyfyMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notyfyMes)
			//用户客户端在线用户map
			model.UpdateUserStatus(notyfyMes)
		case message.ChatMesType:
			OutPutMes(mes)
		default:
			fmt.Print("服务器返回未知信息")

		}

		fmt.Println("mes=", mes)
	}
}
