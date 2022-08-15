package process

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加字段， 用于区分conn
	UserId int
}

//编写通知所有在线用户方法
func (u *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onlineUsers, 然后一个一个发送
	for id, up := range userMagr.onlineUsers {
		if id == userId {
			continue
		}
		//
		up.NotifyMeOnline(userId)
	}
}

func (u *UserProcess) NotifyMeOnline(userId int) (err error) {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json Marshal err=", err)
		return

	}

	mes.Data = string(data)

	// 对mes 序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json Marshal err=", err)
		return

	}

	//发送 创建tf实例
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
	return

}

//编写一个函数serverProcessLogin函数，专门处理登录请求
func (u *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//先从mes中取出 mes.Data, 并直接反序列化LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 1先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2声明一个 LoginResMes， 并完成赋值
	var loginResMes message.LoginResMes

	// 到redis中完成验证

	//1. 使用model.MyUserDao 到redis中验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()

		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}

		// 先测试成功， 然后返回具体的信息
	} else {
		loginResMes.Code = 200
		// 用户登陆成功， 将登录成功的用户放入到userMgr
		// 将登陆成功的ID赋值给u
		u.UserId = loginMes.UserId
		userMagr.AddOnlineUser(u)
		//通知其他用户我上线了
		u.NotifyOthersOnlineUser(loginMes.UserId)
		//将id 放入到loginResMes.UsersId
		for id := range userMagr.onlineUsers { //缺少_
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登陆成功")
	}

	// 3将 loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
	}
	// 4将data 赋值给 resMes
	resMes.Data = string(data)

	// 5对resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
	}

	// 6 发送序列化好的数据 将其封装到writePkg中

	// 使用分层模式（MVC）， 先创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)
	return

}

func (u *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//先从mes中取出 mes.Data, 并直接反序列化registerMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 1先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType

	// 2声明一个 LoginResMes， 并完成赋值
	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
			return
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}

	// 3将 loginResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
	}

	// 4将data 赋值给 resMes
	resMes.Data = string(data)

	// 5对resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
	}

	// 6 发送序列化好的数据 将其封装到writePkg中

	// 使用分层模式（MVC）， 先创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	err = tf.WritePkg(data)
	return

}
