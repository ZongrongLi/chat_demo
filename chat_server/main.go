package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/tiancai110a/chat_demo/errno"
	_ "github.com/tiancai110a/chat_demo/errno"
	proto "github.com/tiancai110a/chat_demo/proto"
)

var mgr *UserMgr

func init() {
	initRedis("localhost:6379", 16, 1024, time.Second*300)
	mgr = NewUserMgr(pool)
}

func readPackage(conn net.Conn) (msg proto.Message, err error) {
	buff := make([]byte, 1024)
	n, err := conn.Read([]byte(buff[0:4]))

	if err != nil {
		fmt.Println("write data  failed")
		return
	}
	packLen := binary.BigEndian.Uint32(buff[0:4])

	n, err = conn.Read([]byte(buff[0:packLen]))

	if err != nil {
		fmt.Println("write data  failed")
		return
	}
	if n != int(packLen) {
		fmt.Println("read data  not finished", n, packLen)
		err = errors.New("read data not fninshed")
		return
	}

	//	fmt.Println("data:", string(buff[0:packLen]))
	msg = proto.Message{}
	err = json.Unmarshal(buff[0:packLen], &msg)
	if err != nil {
		err = errors.New("msg data Unmarshal failed")
		return
	}
	return

}

func writePackage(conn net.Conn, data []byte) {
	buff := make([]byte, 4)
	packLen := uint32(len(data))

	binary.BigEndian.PutUint32(buff[0:4], packLen)
	n, err := conn.Write(buff)
	if err != nil {
		fmt.Println("write data  failed")
		return
	}
	n, err = conn.Write(data)

	if err != nil {
		fmt.Println("write data  failed")
		return
	}

	if n != int(packLen) {
		fmt.Println("write data  not finished", n, packLen)
		err = errors.New("write data not fninshed")
		return
	}

}

func Login(msg proto.Message) (err error) {
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

func Register(msg proto.Message) (err error) {

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

func LoginResp(conn net.Conn, err error) {
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

func process(conn net.Conn) {
	defer conn.Close()
	msg, err := readPackage(conn)
	defer func() {
		LoginResp(conn, err)
	}()

	if err != nil {
		fmt.Println("readPackage: ", err)
		return
	}

	switch msg.Cmd {
	case proto.UserLogin:
		err = Login(msg)
	case proto.UserRegister:
		err = Register(msg)
	default:
		fmt.Println("unkown cmd")
		return
	}

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:10000")
	if err != nil {
		fmt.Println("listen failed, ", err)
		return
	}

	for {
		conn, err := l.Accept()

		go process(conn)
		if err != nil {
			fmt.Println("accept failed, ", err)
			continue
		}
	}
}
