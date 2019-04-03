package controller

import (
	"chat/common/message"
	"chat/common/operationData"
	"encoding/json"
	"net"
)

type ChatMesController struct {

}

//消息群发
func (this ChatMesController)SendGroupMes(mes *message.Message)  {
	//取出得到的内容
	var chatMes message.ChatMes
	json.Unmarshal([]byte(mes.Data),&chatMes)

	//序列化完毕准备发送
	data,_:=json.Marshal(mes)
	for id,userCont := range userMg.onlineUsers{
		if id == chatMes.UserId{
			continue
		}
		this.SendMes(data,userCont.Conn)
	}
}

//发送消息
func(this ChatMesController)SendMes(data []byte,conn net.Conn)  {
	operData := operationData.PackageOperation{
		Conn:conn,
	}
	operData.WritePackage(data)
}