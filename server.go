package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct{
    Ip      string
    Port    int

	//在线用户列表
	OnLineMap	map[string]*User
	mapLock		sync.RWMutex

	//消息广播的channel
	Massage		chan string

}



//创建一个server的接口
func Newserver(ip string, port int) *Server{
        server := &Server{
                Ip: ip,
                Port: port,
		OnLineMap: make(map[string]*User),
		Massage: make(chan string),
        }

        return server
}

//监听Massage信息，一旦有信息就发送出去
func (this *Server) ListenMsg(){
	for{
		msg := <-this.Massage

		//将信息发送给所有在线用户
		this.mapLock.Lock()
		for _,cli := range this.OnLineMap{
				cli.C <-msg
			}
		this.mapLock.Unlock()
	}
}


//广播消息的方法
func (this *Server) BroadCast(user *User, msg string){
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Massage <-sendMsg
}


func (this *Server) Hander(conn net.Conn){
        //用户上线，加入onlinemap
	user := Newuser(conn,this)

	user.OnLine()
	//接收客户端发送的消息

	//监听用户是否活跃
	isLive := make(chan bool)
	go func() {
		buff := make([]byte,4096)
		for {
			n, err := conn.Read(buff)
			if n== 0{
				user.OffLine()
				return
			}

			if err != nil && err != io.EOF{
				fmt.Println("Conn Read Err:",err)
				return
			}

			msg := string(buff[:n-1])
			user.DoMsg(msg)
			isLive <- true
		}
	}()
	

	for {
		select{
		case <- isLive:

		case <-time.After(time.Second*600):

		user.SendMsg("你已超时，已被踢出")
		//user.OffLine()
		close(user.C)
		conn.Close()
		return
		}
	}
}

func (this *Server) Start(){

        listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",this.Ip, this.Port))
        if err != nil{
                fmt.Println("net.Listen err:",err)
                return
        }

        defer listener.Close()

	go this.ListenMsg()

        for{
                conn, err := listener.Accept()
                if err != nil{
                        fmt.Println("LIsteb Accept err:",err)
                        continue
                }
		go this.Hander(conn)
	}
}
