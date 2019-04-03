package controller

import (
	"chat/client/model"
	"chat/common/message"
	"chat/common/operationData"
	"encoding/json"
	"fmt"
)

type ChatMesController struct {

}

//群发信息
func (this ChatMesController)SendGroupMes(content string)  {

	var mes message.Message
	mes.Type = message.ChatMesType

	var chatMes message.ChatMes
	chatMes.Content = content
	chatMes.UserId = model.ConUser.UserId
	chatMes.UserName =model.ConUser.UserName

	//序列化
	data,_:=json.Marshal(chatMes)

	mes.Data = string(data)
	//再次序列化
	data,_=json.Marshal(mes)

	operData := operationData.PackageOperation{
		Conn:model.ConUser.Conn,
	}

	operData.WritePackage(data)

}
//显示收到的消息
func OutPutMes(mes message.Message)  {
	var chatMes message.ChatMes
	json.Unmarshal([]byte(mes.Data),&chatMes)
	//显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s",
		chatMes.UserId, chatMes.Content)
	fmt.Println(info)
	fmt.Println()

}