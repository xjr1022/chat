package message


//用户结构
type User struct {
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
	//用户状态
	UserStatus int `json:"userStatus"`
}





