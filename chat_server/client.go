package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/tiancai110a/chat_demo/errno"
	proto "github.com/tiancai110a/chat_demo/proto"
	"github.com/tiancai110a/chat_demo/transport"
)

type client struct {
	conn net.Conn
}

func (c *client) Login(msg proto.Message) (respdata []byte, err error) {
	cmd := proto.LoginCmd{}
	err = json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		fmt.Println("unmarshal failed: ", err)
		return
	}

	fmt.Println("login: ", cmd)

	u, err := mgr.Login(cmd.Id, cmd.Passwd)

	if err != nil {
		fmt.Println(" mgr.Login failed: ", err)
		return
	}

	clientMgr.AddClient(cmd.Id, c)
	fmt.Println(u)

	lc := proto.LoginCmdRes{}

	usermp := clientMgr.GetAllUsers()
	for id, cli := range usermp {
		if id == cmd.Id {
			continue
		}
		lc.Users = append(lc.Users, id)
		cli.Notify(cmd.Id)
	}
	//user list
	respdata, err = json.Marshal(lc)
	if err != nil {
		fmt.Println("marshal failed: ", err)
		return
	}

	return
}

func (c *client) Notify(userid int) (err error) {
	onlineNotify := proto.UserStatusNotify{Userid: userid, Status: proto.UserOnline}

	msg := proto.Message{}

	data, err := json.Marshal(onlineNotify)
	msg.Cmd = proto.UserNotifyStatus
	msg.Data = data
	transport.SendMessage(c.conn, msg)
	return
}

func (c *client) Register(msg proto.Message) (err error) {

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

func (c *client) sendResp(conn net.Conn, cmd proto.Command, data []byte, err error) {

	msg := proto.Message{}
	msg.Cmd = cmd
	msg.Data = data

	if err != nil {
		if errnum, ok := err.(errno.Errno); !ok {
			msg.Code = -1
		} else {
			msg.Code = errnum.Code
		}

		msg.Error = err.Error()
	} else {
		msg.Code = errno.OK.Code
		msg.Error = errno.OK.Message
	}

	msgdata, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("msg Marshal failed: ", err)
		return
	}
	transport.WritePackage(conn, msgdata)
	return
}
