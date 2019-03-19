package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//定义一个userModel全局变量
var(
	MyUserModel *UserModel
)

//用户结构
type User struct {
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
}

//用户模型
type UserModel struct {
	Pool *redis.Pool
}



//根据id获取用户数据
func (this *UserModel)getUserById(conn redis.Conn ,id int)(user *User,err error)  {
	res,err :=redis.String(conn.Do("HGet","users",id))
	if err != nil {
		if err == redis.ErrNil{
			//表示对应id数据不存在
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	err = json.Unmarshal([]byte(res),user)
	if err != nil {
		fmt.Println("user 序列化失败")
		return
	}
	return
}

//用户登录校验
func(this *UserModel) Login(userId int,userPwd string)(user *User,err error)  {
	//先从连接池取出连接
	conn :=this.Pool.Get()
	defer conn.Close()
	user,err = this.getUserById(conn,userId)
	if err != nil {
		return
	}
	if user.UserPwd != userPwd{
		err = ERROR_USER_PWD

	}
	return
}

