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
	fmt.Scanf("%d", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
	case 2:
		fmt.Println("发送消息")
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
		default:
			fmt.Print("服务器返回未知信息")

		}

		fmt.Println("mes=", mes)
	}
}
