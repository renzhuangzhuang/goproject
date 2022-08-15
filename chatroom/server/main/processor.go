package main

import (
	"chatroom/common/message"
	"chatroom/server/process"
	"chatroom/server/utils"
	"errors"
	"fmt"
	"net"
)

// 先创建一个Process的结构体
type Process struct {
	Conn net.Conn
}

// 功能 根据客户端发送消息种类不同， 决定调用哪个函数
func (p *Process) ServerProcessMes(mes *message.Message) (err error) {

	//看看能否接收到客户端发送的群发消息
	fmt.Println("mes=", mes)
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//实例一个UserProcess
		up := &process.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &process.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//完成转发消息的功能
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (p *Process) process2() error {

	//循环读取客户端发送的信息
	for {
		// 这里读取数据包封装成一个函数readPkg()
		// 创建Transfer
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			err = errors.New("read pkg header error")
			fmt.Println(err)
			return err
		}
		err = p.ServerProcessMes(&mes)
		if err != nil {
			return err
		}

	}
}
