package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"chat/common/message"
)

//用户模型
type UserModel struct {
	Pool *redis.Pool
}
//定义一个userModel全局变量
var(
	MyUserModel *UserModel
)


//根据id获取用户数据
func (this *UserModel)getUserById(conn redis.Conn ,id int)(user *message.User,err error)  {
	res,err :=redis.String(conn.Do("HGet","users",id))
	if err != nil {
		if err == redis.ErrNil{
			//表示对应id数据不存在
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &message.User{}
	err = json.Unmarshal([]byte(res),user)
	if err != nil {
		fmt.Println("user 序列化失败")
		return
	}
	return
}

//用户登录校验
func(this *UserModel) Login(userId int,userPwd string)(user *message.User,err error)  {
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

//用户注册
func(this *UserModel) Register(user *message.User)(err error)  {
	//先从连接池取出连接
	conn :=this.Pool.Get()
	defer conn.Close()
	_,err = this.getUserById(conn,user.UserId)
	//如果不为空，说明用户已经存在
	if err == nil {
		fmt.Println("注册前数据库查询是否重复 err=",err)
		return
	}
	//完成注册
	data,err :=json.Marshal(user)
	_,err = conn.Do("Hset","users",user.UserId,string(data))
	if err != nil {
		fmt.Println("redis保存注册信息失败 err=",err)
		return
	}
	return
}