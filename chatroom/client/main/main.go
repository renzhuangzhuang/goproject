package main

import (
	"chatroom/client/process"
	"fmt"
	"os"
)

// 定义两个变量， 一个id 一个密码
var userId int
var userPwd string
var userName string

func main() {
	// 接收用户的选择
	var key int
	//判断是否还继续选择菜单
	//var loop = true

	for {
		fmt.Println("---------------欢迎登陆多人聊天室-------------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			//完成登陆
			//创建userProcess实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名称")
			fmt.Scanf("%s\n", &userName)
			// 调用UserProcess完成实例
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
			//loop = false
		default:
			fmt.Println("你的输入有误 请重新输入")

		}
	}

}
