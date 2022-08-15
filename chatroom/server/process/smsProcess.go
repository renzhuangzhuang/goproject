package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

func (s *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的map
	//将消息转发出去
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	for id, up := range userMagr.onlineUsers {
		//过滤自己消息
		if id == smsMes.UserId {
			continue
		}
		s.SendMesTOEachOnlineUser(data, up.Conn)
	}

}

func (s *SmsProcess) SendMesTOEachOnlineUser(data []byte, conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err =", err)
	}
}
