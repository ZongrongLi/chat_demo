package transport

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	proto "github.com/tiancai110a/chat_demo/proto"
)

func ReadPackage(conn net.Conn) (msg proto.Message, err error) {
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

func WritePackage(conn net.Conn, data []byte) {
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

func SendMessage(conn net.Conn, msg proto.Message) {
	msgdata, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("msg Marshal failed: ", err)
		return
	}
	WritePackage(conn, msgdata)
}
