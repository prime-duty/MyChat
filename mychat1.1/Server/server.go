package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string           //ip地址
	Port      int              //端口号
	Massage   chan string      //消息通道
	OnlineMap map[string]*User //在线用户列表
	maplock   sync.RWMutex
}

//Server的构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		Massage:   make(chan string),
		OnlineMap: make(map[string]*User),
	}

	return server
}

//监听Massage，一有消息就发送出去
func (this *Server) ListenMsg() {
	for {
		msg := <-this.Massage
		this.maplock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.maplock.Unlock()
	}
}

func (this *Server) Hander(conn net.Conn) {
	user := NewUser(conn, this)

	go func() {
		buff := make([]byte, 4096)

		for {
			n, err := conn.Read(buff)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil {
				return
			}
			msg := string(buff[:n-1])
			user.DoMsg(msg)
		}
	}()

}

func (this *Server) SendMsgAll(user *User, msg string) {
	smsg := "[" + user.Name + "]" + ":" + msg
	this.Massage <- smsg
}

func (this *Server) Start() {
	listenfd, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("Listen err")
		return
	}

	defer listenfd.Close()

	go this.ListenMsg()

	for {
		conn, err := listenfd.Accept()
		if err != nil {
			fmt.Println("Accept err")
			continue
		}

		go this.Hander(conn)
	}

}
