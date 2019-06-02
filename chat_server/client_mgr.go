package main

import "fmt"

type ClientMgr struct {
	onlineUsers map[int]*client
}

var clientMgr ClientMgr

func init() {
	clientMgr = ClientMgr{}
	clientMgr.onlineUsers = make(map[int]*client)
}

func (c *ClientMgr) AddClient(userid int, cli *client) {
	c.onlineUsers[userid] = cli
}

func (c *ClientMgr) GetClient(userid int) (cli *client, err error) {
	cli, ok := c.onlineUsers[userid]
	if !ok {
		err = fmt.Errorf("user %d not exist", userid)
		return
	}

	return
}

func (c *ClientMgr) GetAllUsers() map[int]*client {
	return c.onlineUsers
}

func (c *ClientMgr) DeleteClient(userid int) {
	delete(c.onlineUsers, userid)
}
