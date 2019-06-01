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
	fmt.Println(11)
	n, err := conn.Read([]byte(buff[0:4]))

	if err != nil {
		fmt.Println("write data  failed")
		return
	}
	fmt.Println(12)
	packLen := binary.BigEndian.Uint32(buff[0:4])

	n, err = conn.Read([]byte(buff[0:packLen]))

	if err != nil {
		fmt.Println("write data  failed")
		return
	}
	fmt.Println(13)
	if n != int(packLen) {
		fmt.Println("read data  not finished", n, packLen)
		err = errors.New("read data not fninshed")
		return
	}
	fmt.Println(14)

	//	fmt.Println("data:", string(buff[0:packLen]))
	msg = proto.Message{}
	err = json.Unmarshal(buff[0:packLen], &msg)
	if err != nil {
		err = errors.New("msg data Unmarshal failed")
		return
	}
	fmt.Println(15)
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

func Login(msg proto.Message) {

	fmt.Println(2)
	cmd := proto.LoginCmd{}
	err := json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		fmt.Println("unmarshal failed: ", msg.Data)
		return
	}
	fmt.Println(3)

	fmt.Println("login: ", cmd)
}

func Register(msg proto.Message) {

	reg := proto.RegisterCmd{}
	err := json.Unmarshal([]byte(msg.Data), &reg)
	if err != nil {
		fmt.Println("unmarshal failed: ", msg.Data)

		return
	}

	fmt.Println("login: ", reg)
}

func LoginResp(conn net.Conn, code int, Error string) (err error) {
	lc := proto.LoginCmdRes{}
	lc.Code = code
	lc.Error = Error
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
	if err != nil {
		fmt.Println("readPackage: ", err)
		return
	}
	Login(msg)
	err = LoginResp(conn, 10, "no error")
	if err != nil {
		fmt.Println("Error LoginResp", err.Error())
		return
	}
	//writePackage(conn, []byte("hello world"))

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
