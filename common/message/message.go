package message

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

//定义消息结构体
type LoginMes struct {
	UserId   int    `json:"userid"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}
type LoginResMes struct {
	Code  int    `json:"code"`  //返回状态码
	Error string `json:"error"` //返回错误信息
}
