package main

import "net"
import "strings"

type User struct{
	Name string
	Addr string
	C    chan string
	conn net.Conn
	server *Server
}


func Newuser(conn net.Conn, server *Server) *User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}
	
	go user.ListenMassage()
	return user

}

//用户上线功能
func (this *User)OnLine(){

        this.server.mapLock.Lock()
        this.server.OnLineMap[this.Name] = this
        this.server.mapLock.Unlock()

        //广播用户上线信息
        this.server.BroadCast(this,"已上线")
}

//用户下线
func (this *User) OffLine(){
	this.server.mapLock.Lock()
	delete(this.server.OnLineMap,this.Name)
        this.server.mapLock.Unlock()

        //广播用户下线信息
        this.server.BroadCast(this,"已下线")

}

//给当前用户发消息
func (this *User)SendMsg(msg string){
	this.conn.Write([]byte(msg))
}


//用户处理消息
func (this *User)DoMsg(msg string){
	if msg == "who"{
		this.server.mapLock.Lock()
		for _,user := range this.server.OnLineMap{
			onlinemsg := "[" + user.Addr + "]" + user.Name + ":" + "Online\n"
			this.SendMsg(onlinemsg)
		}
		this.server.mapLock.Unlock()
	}else if len(msg) > 7 && msg[:7] == "rename|"{
		newName := strings.Split(msg,"|")[1]
		_, ok := this.server.OnLineMap[newName]
		if ok {
			this.SendMsg("此用户名已被占用")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnLineMap,this.Name)
			this.server.OnLineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("更新用户名成功\n")
		}
	}else if len(msg) > 4 && msg[:3] == "to|"{
		rName := strings.Split(msg,"|")[1]
		if rName == ""{
			this.SendMsg("请使用正确的命令，如\"to|张三|nihao\"。\n")
			return
		}

		rUser, ok := this.server.OnLineMap[rName]
		if !ok {
			this.SendMsg("该用户名不存在")
			return
		}

		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("发送内容不能为空")
			return
		}
		rUser.SendMsg(this.Name + "对您说：" + content+"\n")
	}else{
		this.server.BroadCast(this,msg)
	}
}

func (this *User) ListenMassage(){
	for{
		msg := <-this.C
		this.conn.Write([]byte(msg+"\n"))
	}

}
