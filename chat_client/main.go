package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	proto "github.com/tiancai110a/chat_demo/proto"
)

func readPackage(conn net.Conn) (msg proto.Message, err error) {
	buff := make([]byte, 100)

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
		fmt.Println("write data  not finished")
		err = errors.New("write data not fninshed")
		return
	}

}

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

	msgdata, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("msg Marshal failed: ", err)
		return
	}
	writePackage(conn, msgdata)
	return

}

func Register(conn net.Conn, id int, passwd string) (err error) {

	reg := proto.RegisterCmd{}
	reg.User = 10

	data, err := json.Marshal(reg)
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

func GetResp(conn net.Conn) {
	msg, err := readPackage(conn)
	fmt.Println("data len: ", len(msg.Data))

	if err != nil {
		fmt.Println("login read data  failed")
		return
	}
	lc := proto.LoginCmdRes{}
	err = json.Unmarshal([]byte(msg.Data), &lc)
	if err != nil {
		fmt.Println("unmarshal failed: ", lc)
		return
	}

	fmt.Println("resp: ", lc)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:10000")

	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	err = Register(conn, 2, "passwd")
	if err != nil {
		fmt.Println("Error Login", err.Error())
		return
	}

	GetResp(conn)

}
