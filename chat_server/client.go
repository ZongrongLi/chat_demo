package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/tiancai110a/chat_demo/errno"
	proto "github.com/tiancai110a/chat_demo/proto"
)

type client struct {
}

func (c client) Login(msg proto.Message) (err error) {
	cmd := proto.LoginCmd{}
	err = json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		fmt.Println("unmarshal failed: ", msg.Data)
		return
	}

	fmt.Println("login: ", cmd)

	u, err := mgr.Login(cmd.Id, cmd.Passwd)
	fmt.Println(u)
	return
}

func (c client) Register(msg proto.Message) (err error) {

	reg := proto.RegisterCmd{}
	err = json.Unmarshal([]byte(msg.Data), &reg)
	if err != nil {
		fmt.Println("unmarshal failed: ", err)
		return
	}
	fmt.Println("regist data", reg.User)
	err = mgr.Register(&reg.User)
	fmt.Println("register: ", reg)
	return
}

func (c client) LoginResp(conn net.Conn, err error) {
	lc := proto.LoginCmdRes{}
	if err != nil {
		if errnum, ok := err.(errno.Errno); !ok {
			lc.Code = -1
		} else {
			lc.Code = errnum.Code
		}

		lc.Error = err.Error()
	} else {
		lc.Code = errno.OK.Code
		lc.Error = errno.OK.Message
	}
	data, err := json.Marshal(lc)
	if err != nil {
		return
	}
	fmt.Println("data: ", data)

	msg := proto.Message{}
	msg.Cmd = proto.UserRegister
	msg.Data = data

	msgdata, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("msg Marshal failed: ", err)
		return
	}
	writePackage(conn, msgdata)
	return
}
