package process

import (
	"fmt"
)

//因为UserMgr实例在服务端有且只有一个
//将其定义为全局变量
var (
	userMagr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

func init() {
	userMagr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUser添加
func (u *UserMgr) AddOnlineUser(up *UserProcess) {
	u.onlineUsers[up.UserId] = up
}

//删除
func (u *UserMgr) DeleteOnlineUser(userId int) {
	delete(u.onlineUsers, userId)
}

//返回当前所以在线用户
func (u *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return u.onlineUsers
}

//根据Id返回对应的值
func (u *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//
	up, ok := u.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
