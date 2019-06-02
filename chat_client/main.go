package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	model "github.com/tiancai110a/chat_demo/model"
	proto "github.com/tiancai110a/chat_demo/proto"
	"github.com/tiancai110a/chat_demo/transport"
)

func Login(conn net.Conn, id int, passwd string) (err error) {

	c := proto.LoginCmd{}
	c.Id = id
	c.Passwd = passwd
	data, err := json.Marshal(c)
	if err != nil {
		fmt.Println("login data Marshal failed: ", err)

		return
	}

	msg := proto.Message{}
	msg.Cmd = proto.UserLogin
	msg.Data = data

	transport.SendMessage(conn, msg)
	return

}

func Register(conn net.Conn, user model.User) (err error) {

	reg := proto.RegisterCmd{}
	reg.User = user

	data, err := json.Marshal(reg)
	if err != nil {
		return
	}
	fmt.Println("data: ", data)

	msg := proto.Message{}
	msg.Cmd = proto.UserRegister
	msg.Data = data

	transport.SendMessage(conn, msg)
	return
}

func GetList(conn net.Conn, userid int) (err error) {

	req := proto.UserListReq{UserId: userid}

	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	fmt.Println("data: ", data)

	msg := proto.Message{}
	msg.Cmd = proto.UserRegister
	msg.Data = data

	transport.SendMessage(conn, msg)
	return
}

func GetResp(conn net.Conn) {

	var err error
begin:
	for err == nil {
		msg, err := transport.ReadPackage(conn)
		fmt.Println("data len: ", len(msg.Data))

		if err != nil {
			fmt.Println("login read data  failed")
			break begin
		}

		switch msg.Cmd {
		case proto.UserLogin:
		case proto.UserLoginRes:
			lc := proto.LoginCmdRes{}
			err = json.Unmarshal([]byte(msg.Data), &lc)
			if err != nil {
				fmt.Println("unmarshal failed: ", lc)
				break begin
			}
			fmt.Println("resp: ", lc)
		case proto.UserRegister:
		case proto.UserNotifyStatus:
			lc := proto.UserStatusNotify{}
			err = json.Unmarshal([]byte(msg.Data), &lc)
			if err != nil {
				fmt.Println("unmarshal failed: ", lc)
				break begin
			}
			fmt.Println("notify: ", lc)
		case proto.DefaultRes:
			fmt.Println("Register: ", msg)

		default:
			break
		}
	}

}

func main() {
	var userid int
	var passswd string
	flag.IntVar(&userid, "u", 1, "please input conf path")
	flag.StringVar(&passswd, "p", "", "please input log level")

	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:10000")

	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}

	// user := model.User{UserId: 2, Passwd: "12345678"}

	// err = Register(conn, user)
	err = Login(conn, userid, passswd)
	if err != nil {
		fmt.Println("Error Login", err.Error())
		return
	}

	go func() {
		GetResp(conn)
	}()

	time.Sleep(1000 * time.Second)
}
