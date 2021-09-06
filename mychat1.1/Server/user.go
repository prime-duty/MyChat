package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis"
)

type User struct {
	Id      string
	Name    string
	C       chan string
	Conn    net.Conn
	Server  *Server
	logined bool
	rdb     *redis.Client
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	err, redb := initClient()
	if err != nil {
		fmt.Println("Redis 异常")
	}
	user := &User{
		Id:      addr,
		Name:    addr,
		C:       make(chan string),
		Conn:    conn,
		Server:  server,
		logined: false,
		rdb:     redb,
	}

	fmt.Println("一个用户上线")
	go user.ListenC()

	return user
}

//Redis操作

// 初始化连接
func initClient() (err error, r *redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err, nil
	}
	return nil, rdb
}

func (this *User) SetVal(id, pwd string) bool {
	err := this.rdb.Set(id, pwd, 0).Err()
	if err != nil {
		return false
	}
	return true
}

func (this *User) GetVal(id string) string {
	val, err := this.rdb.Get(id).Result()
	if err != nil {
		fmt.Println("查询异常", err)
		return ""
	}

	return val
}

//处理用户消息消息
func (this *User) DoMsg(msg string) {
	if msg == "login" {
		fmt.Println("登录状态")
		buff := make([]byte, 4096)
		for {
			n, err := this.Conn.Read(buff)

			if err != nil {
				return
			}
			msg := string(buff[:n-1])
			fmt.Println(msg)
			id := strings.Split(msg, "@")[0]
			pwd := strings.Split(msg, "@")[1]
			if this.GetVal(id) == pwd {
				this.Online(id)
				this.SendMsg("true")

				this.logined = true
				break
			} else {
				this.SendMsg("false")
				continue
			}
		}
	} else if len(msg) >= 3 && msg[:3] == "re|" {
		//this.SendMsg("请设定账号和密码")
		//buff := make([]byte, 4096)

		//n, err := this.Conn.Read(buff)

		// if err != nil {
		// 	return
		// }
		//msg := string(buff[:n-1])
		id := strings.Split(msg, "|")[1]
		pwd := strings.Split(msg, "|")[2]
		fmt.Println(id)
		fmt.Println(pwd)
		if this.GetVal(id) != "" {
			this.SendMsg("err")
		} else {
			if this.SetVal(id, pwd) == true {
				this.SendMsg("true")

			} else {
				this.SendMsg("false")

			}
		}
		//this.SendMsg("false")

	} else if this.logined != true {
		this.SendMsg("You hava not logged,please sign in and try again")
	} else if len(msg) >= 5 && msg[:3] == "set" {
		//fmt.Println("------")
		//this.SendMsg("请输入你想设定的昵称")
		// buff := make([]byte, 4096)
		// for {
		// 	n, err := this.Conn.Read(buff)

		// 	if err != nil {
		// 		return
		// 	}
		// 	msg := string(buff[:n-1])
		// 	delete(this.Server.OnlineMap, this.Name)
		// 	this.Server.OnlineMap[msg] = this
		// 	this.Name = msg
		// 	if this.Server.OnlineMap[msg] == this {
		// 		this.SendMsg("你的昵称已设定成功")
		// 		break
		// 	} else {
		// 		this.SendMsg("你的昵称设定失败，请重新设置")
		// 		continue
		// 	}
		// }
		rname := strings.Split(msg, "|")[1]
		this.Server.maplock.Lock()
		delete(this.Server.OnlineMap, this.Name)
		this.Server.OnlineMap[rname] = this
		this.Server.maplock.Unlock()
		this.Name = rname
	} else if msg == "ls" {
		this.Server.maplock.Lock()
		//this.SendMsg("在线人员昵称如下：")
		for name, _ := range this.Server.OnlineMap {
			if name == this.Name {
				continue
			}
			this.SendMsg(name)
		}
		this.Server.maplock.Unlock()
		this.SendMsg("-------------------------------")
	} else if len(msg) > 3 && msg[:3] == "to|" {
		name := strings.Split(msg, "|")[1]
		sendmsg := strings.Split(msg, "|")[2]
		senduser := this.Server.OnlineMap[name]
		if senduser == nil {
			this.SendMsg("用户昵称不存在")
		} else {
			senduser.SendMsg("[" + this.Name + "私聊你：" + "]" + sendmsg)
		}
	} else {
		this.Server.SendMsgAll(this, msg)
	}
}

//用户上线
func (this *User) Online(name string) {
	this.Name = name
	this.Server.maplock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.maplock.Unlock()

	this.Server.SendMsgAll(this, "已上线")
}

//用户离线
func (this *User) Offline() {
	this.Server.maplock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.maplock.Unlock()

	this.Server.SendMsgAll(this, "已离线")
}

//监听消息
func (this *User) ListenC() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}

//给当前用户发消息
func (this *User) SendMsg(msg string) {
	smsg := msg + "\n"
	this.Conn.Write([]byte(smsg))
}
