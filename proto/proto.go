package proto

import "github.com/tiancai110a/chat_demo/model"

type Command byte

const (
	DefaultRes Command = iota
	UserLogin
	UserLoginRes
	UserRegister
	UserNotifyStatus
)

type UserStatus byte

const (
	UserOnline = iota
	UserOffline
)

type Base struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
type Message struct {
	Cmd  Command `json:"cmd"`
	Data []byte  `json:"data"`
	Base
}

type LoginCmd struct {
	Id     int    `json:"user_id"`
	Passwd string `json:"passwd"`
}

type RegisterCmd struct {
	User model.User `json:"user"`
}

type LoginCmdRes struct {
	Users []int `json:"users"`
}

type UserStatusNotify struct {
	Userid int        `json:"userid"`
	Status UserStatus `json:"status"`
}

type UserListReq struct {
	UserId int `json:"userid"`
}
