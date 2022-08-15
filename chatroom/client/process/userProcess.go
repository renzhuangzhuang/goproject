package process

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要字段
}

func (u *UserProcess) Login(userId int, userPwd string) (err error) {

	//下一个就要开始定协议
	/* fmt.Printf(" userId = %d userPwd= %s\n", userId, userPwd)
	return nil */
	//1.链接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return

	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn 发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//4将结构体loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 5 把data赋值给 mes.Data字段
	mes.Data = string(data)

	//6 将mes进行序列化化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 7 data为需要发送的数据 先发送一个长度 再发送数据
	// 7.1 先把data的长度发送给服务器
	// 先获取data长度，再转成一个表示长度的byte切片
	pkgLen := uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], pkgLen)
	_, err = conn.Write(buf[:])
	if err != nil {
		fmt.Println("conn.Write err = ", err)
		return
	}
	//fmt.Printf("客户端，发送消息的长度=%d 内容=%s", len(data), string(data))
	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write err = ", err)
		return
	}

	//这里需要处理服务器端返回的消息
	// 实例化一个Transfer
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg() //
	if err != nil {
		fmt.Println("readPkg(conn) err = ", err)
		return
	}
	//mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//fmt.Println("登陆成功")
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//显示登陆成功后菜单 循环显示
		// 可以显示在线用户列表
		fmt.Println("当前在线用户列表：")
		for _, v := range loginResMes.UsersId {
			// 不显示自己
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成 客户端的onlineUsers 完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		// 启用一个协程和服务器 保持通信，如果服务器有数据推送给客户端
		// 则接收并显示在客户端

		go serverProcessMes(conn)
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

// 完成请求注册的方法
func (u *UserProcess) Register(userId int,
	userPwd string, userName string) (err error) {
	//1.链接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return

	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn 发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建一个LoginMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	//4将结构体RegisterMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5 把data赋值给 mes.Data字段
	mes.Data = string(data)
	//fmt.Println(data)
	//6 将mes进行序列化化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//创建实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	mes, err = tf.ReadPkg() //
	if err != nil {
		fmt.Println("readPkg(conn) err = ", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)

	if registerResMes.Code == 200 {
		fmt.Println("注册成功,你重新登录一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return

}
