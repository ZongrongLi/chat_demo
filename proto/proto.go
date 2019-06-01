package proto

type command byte

const (
	UserLogin command = iota
	UserLoginRes
	UserRegister
)

type Message struct {
	Cmd  command `json:"cmd"`
	Data []byte  `json:"data"`
}

type LoginCmd struct {
	Id     int    `json:"user_id"`
	Passwd string `json:"passwd"`
}

type RegisterCmd struct {
	User int `json:"user"`
}

type LoginCmdRes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
