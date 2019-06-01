package main

import (
	"fmt"
	"net"

	proto "github.com/tiancai110a/chat_demo/proto"
)

func process(conn net.Conn) {
	cli := client{}
	defer conn.Close()
	msg, err := readPackage(conn)
	defer func() {
		cli.LoginResp(conn, err)
	}()

	if err != nil {
		fmt.Println("readPackage: ", err)
		return
	}

	switch msg.Cmd {
	case proto.UserLogin:
		err = cli.Login(msg)
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
