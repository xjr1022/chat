package main

import (
	"fmt"
	"io"
	"net"
	"project/chatRoom/common/message"
	"project/chatRoom/common/operationData"
	"project/chatRoom/server/controller"
)

type Processor struct {
	Conn net.Conn
}



func (this *Processor)MainProcess(conn net.Conn) (err error){

	//循环处理客户端发送的消息
	for {
		packOper := &operationData.PackageOperation{
			Conn:conn,
		}
		 mes , err := packOper.ReadPackage()

		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端退出")
			}else{
				fmt.Println("conn.read err=", err)
			}
			return err
		}
		fmt.Println("读取到的mes=", mes)
		err=this.ServerProcessMes(&mes)
		if err != nil {
			return err
		}

	}

}

//给接收的消息分配相应的函数处理
func (this *Processor)ServerProcessMes(mes *message.Message)(err error){
	userController:=controller.UserController{
		Conn:this.Conn,
	}
	switch mes.Type {
	case message.LoginMesType:
		//处理登录

		err =   userController.ServerProcessLogin(mes)
	case message.RegisterMesType:
		err = userController.ServerProcessRegister(mes)
	default:
		fmt.Println("消息类型不存在")
	}
	return
}