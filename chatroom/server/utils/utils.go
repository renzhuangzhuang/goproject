package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 这里将方法关联到结构体中
type Transfer struct {
	// 分析需要哪些字段
	Conn net.Conn
	Buf  [8096]byte //这是传输时候，使用的缓存

}

func (t *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送数据")
	_, err = t.Conn.Read(t.Buf[:4])
	if err != nil {
		fmt.Println("conn.Read err=", err)
		return
	}

	//根据buf[:4] 转成一个uint32
	pkglen := binary.BigEndian.Uint32(t.Buf[:4])

	//根据pkglen读取消息内容
	n, err := t.Conn.Read(t.Buf[:pkglen])
	if n != int(pkglen) || err != nil {
		fmt.Println("conn.Read err=", err)
	}

	//反序列化pkglen
	//先反序列化成Message
	err = json.Unmarshal(t.Buf[:pkglen], &mes) //mes要加取&
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}
	return

}

//
func (t *Transfer) WritePkg(data []byte) (err error) {

	// 1 先发送一个长度给对方，
	pkgLen := uint32(len(data))
	binary.BigEndian.PutUint32(t.Buf[0:4], pkgLen)
	_, err = t.Conn.Write(t.Buf[:4])
	if err != nil {
		fmt.Println("conn.Write err = ", err)
		return
	}
	//发送data本身
	n, err := t.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write err = ", err)
		return
	}
	return
}
