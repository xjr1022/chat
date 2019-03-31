package controller

import "fmt"

/*
这个控制器的作用主要是对onlineUser的增删改查
 */

 /*
 userManager实例在服务器端只有一个，因此将其定义为全局变量
  */
var (
	userMg *UserManager
)

//保存在线用户的标识
type UserManager struct {
	onlineUsers map[int]*UserController
}

//对userMg初始化
func init()  {
	userMg = &UserManager{
		onlineUsers:make(map[int] *UserController,1024),
	}
}
//修改和添加
func (this *UserManager)AddOnlineUser(userController *UserController)  {
	this.onlineUsers[userController.UserId] = userController
}

//删除
func (this *UserManager)DelOnlineUser(userId int)  {
	delete(this.onlineUsers,userId)
}

//返回当前所有在线的用户
func (this *UserManager)GetAllOnlineUser()map[int]*UserController  {
	return this.onlineUsers
}

//根据id返回对应的用户conn
func (this *UserManager)GetOnlineUserById(userId int)(userController *UserController,err error)  {
	//从map中取值，并判断是否成功取到
	userController,ok:=this.onlineUsers[userId]
	//如果为false，则说明用户不在线
	if !ok {
		err = fmt.Errorf("用户 %d 不在线",userId)
		return
	}
	return
}