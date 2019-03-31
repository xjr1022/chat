package message

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

//定义登录消息结构体
type LoginMes struct {
	UserId   int    `json:"userid"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}
//定义登录消息响应结构体
type LoginResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //返回错误信息
	OnlineUsersId []int `json:"onlineUsersId"`//在线用户id
}
//定义注册结构体
type RegisterMes struct {
	User User `json:"user"`
}
//定义注册响应结构体
type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //返回错误信息
}
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` 
	Status int `json:"status"`
} 