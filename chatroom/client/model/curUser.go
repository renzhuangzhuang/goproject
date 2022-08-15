package model

import (
	"chatroom/common/message"
	"net"
)

// 将CurUser 做成全局变量

type CurUser struct {
	Conn net.Conn
	message.User
}
