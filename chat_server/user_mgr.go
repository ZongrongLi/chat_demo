package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/tiancai110a/chat_demo/errno"
	"github.com/tiancai110a/chat_demo/model"
)

var UserTable = "users"

type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (mgr *UserMgr) {
	mgr = &UserMgr{
		pool: pool,
	}
	return
}

func (u *UserMgr) GetUser(id int) (user *model.User, err error) {
	conn := u.pool.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("HGet", UserTable, fmt.Sprintf("%d", id)))

	if err != nil {
		if err == redis.ErrNil {
			err = errno.ErrUserNotExist
			return
		}
		fmt.Println("getuser from redis failed.err:", err)
		return
	}
	user = &model.User{}
	err = json.Unmarshal([]byte(result), user)

	return

}

func (u *UserMgr) Login(id int, passwd string) (user *model.User, err error) {

	user, err = u.GetUser(id)
	if err != nil {
		return
	}

	if user.UserId != id || user.Passwd != passwd {
		fmt.Println("login info wrong")
		err = errno.ErrInvalidPasswd
	}

	user.Status = model.UserStatusOnline
	user.LastLogin = fmt.Sprintf("%v", time.Now())

	return

}

func (u *UserMgr) Register(user *model.User) (err error) {

	_, err = u.GetUser(user.UserId)
	if err == nil {
		err = errno.ErrUserExist
		fmt.Println("id already exist")
		return
	}

	if err != errno.ErrUserNotExist {
		fmt.Println("register err: ", err)
		return
	}

	data, err := json.Marshal(user)

	if err != nil {
		fmt.Println("user marshal failed", err)
		return
	}
	conn := u.pool.Get()
	defer conn.Close()
	_, err = redis.String(conn.Do("HSet", UserTable, fmt.Sprintf("%d", user.UserId), string(data)))
	return
}
