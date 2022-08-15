package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// 在服务器启动后， 就初始化一个UserDao 做成全局变量, 在需要和redis操作时 就直接使用
var (
	MyUserDao *UserDao
)

//定义一个UserDao 结构体
//完成对User 结构体的各自操作
type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式 创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return

}

//1.根据一个用户Id 返回一个User实例
func (u *UserDao) GetUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过给定ID 去redis查询用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			//表示users哈希中， 没有找到对应Id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	// 将res反序列化
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}
	return

}

//完成登陆校验 Login
//1. Login 完成对用户的验证
//2. 如果ID 和PWD 正确，返回一个user实例
//3. 错误返回错误信息

func (u *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 先从UserDao 的连接池中取出一根连接
	conn := u.pool.Get()
	defer conn.Close()
	user, err = u.GetUserById(conn, userId)
	if err != nil {
		return
	}
	// 证明这个用户获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}

	return

}
func (u *UserDao) Register(user *message.User) (err error) {
	// 先从UserDao 的连接池中取出一根连接
	conn := u.pool.Get()
	defer conn.Close()
	_, err = u.GetUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	// 这时说明 id 在redis还没有，则可以完成注册
	data, err := json.Marshal(user) // 序列化
	if err != nil {
		return
	}
	//可以入库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err = ", err)
		return
	}
	return

}
