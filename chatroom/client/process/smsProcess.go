package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊
func (s *SmsProcess) SendGroupMes(content string) (err error) {

	// 1 创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserId

	//3、序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal err =  ", err.Error())

		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =  ", err)
		return
	}

	// 将 mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送消息有问题")
		return

	}
	return

}
