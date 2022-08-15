package process

import (
	"chatroom/client/model"
	"chatroom/common/message"
	"fmt"
)

// 客户端维护的map

var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

var CurUser model.CurUser //在用户登录成功后完成初始化

// 在客户端显示online用户
func outputOnlineUser() {
	fmt.Println("当前在线用户列表")
	for id := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

//编写一个方法 处理返回的NotifyUserStatusMes
func updataUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

}
