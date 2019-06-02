package main

import (
	"fmt"
	"net"

	proto "github.com/tiancai110a/chat_demo/proto"
	"github.com/tiancai110a/chat_demo/transport"
)

func process(conn net.Conn) {
	cli := client{conn: conn}
	//defer conn.Close()
	var data []byte

	msg, err := transport.ReadPackage(conn)
	command := proto.DefaultRes
	defer func() {
		cli.sendResp(conn, command, data, err)
	}()

	if err != nil {
		fmt.Println("readPackage: ", err)

		conn.Close()
		return
	}
	switch msg.Cmd {
	case proto.UserLogin:
		command = proto.UserLoginRes
		data, err = cli.Login(msg)
	case proto.UserRegister:
		err = cli.Register(msg)
	default:
		fmt.Println("unkown cmd")
		return
	}

}
func runServer(addr string) (err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen failed, ", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept failed, ", err)
			continue
		}

		go process(conn)
	}
}
