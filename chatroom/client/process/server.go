package process

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登陆成功后的界面..
func ShowMenu() {
	fmt.Println("------------恭喜XXX登陆成功------------")
	fmt.Println("------------1. 显示在线用户列表------------")
	fmt.Println("------------2. 发送消息------------")
	fmt.Println("------------3. 信息列表------------")
	fmt.Println("------------4. 退出系统------------")
	fmt.Println("请选择(1-4):")
	var key int
	var content string
	fmt.Scanf("%d\n", &key)
	smsProcess := &SmsProcess{}
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("请输入你想对大家说的话")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")

	case 4:
		fmt.Println("你选择了退出系统")
		os.Exit(0)
	default:
		fmt.Println("你输入的有误")
	}
}

func serverProcessMes(conn net.Conn) {
	// 和服务器端保持通信
	// 创建一个transfer实例，不听读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err = ", err)
			return
		}

		//fmt.Printf("mes = %v\n", mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 处理
			//1. 取出NotifyUserStatusMesType
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2. 把这个用户，保存到客户map[int]User中
			updataUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			//有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器返回来未知的消息类型")

		}
	}

}
