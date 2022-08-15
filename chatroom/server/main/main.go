package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

//处理和客户端的通信
func process1(conn net.Conn) (err error) {
	// 需要延时关闭
	defer conn.Close()
	// 调用总控
	process := &Process{
		Conn: conn,
	}
	err = process.process2()
	if err != nil {
		fmt.Println("客户端和服务器端协程出错", err)
		return
	}
	return
}

// 编写一个函数初始化UserDao
func initUserDao() {
	// 这里pool 就是一个全局变量
	// 注意初始化顺序问题

	model.MyUserDao = model.NewUserDao(pool)
}
func main() {
	//提示信息
	// 当服务器启动时 就初始化redis连接池
	initPool("0.0.0.0:6379", 16, 0, 300*time.Second)
	initUserDao()
	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	//defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	for {
		fmt.Println("等待客户端链接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦链接成功， 则启用一个协程和客户端保持通信
		go process1(conn)
	}

}
