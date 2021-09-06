package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Cli struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
	login      bool
}

//cli的构造函数
func NewCli(ip string, port int) *Cli {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Conn err")
		return nil
	}
	cli := &Cli{
		ServerIp:   ip,
		ServerPort: port,
		conn:       conn,
		flag:       0,
		login:      false,
	}

	return cli
}

func (this *Cli) menu() bool {
	var flag int
	if this.login == false {
		fmt.Println("你还没登录，请先登录")
		fmt.Println("-------------------------------")
		fmt.Println("-    登录      -      0       -")
		fmt.Println("-    注册      -      1       -")
		fmt.Println("-------------------------------")
		fmt.Print("请选择：")
		fmt.Scan(&flag)
		if flag == 0 || flag == 1 {
			this.flag = flag
			return true
		} else {
			fmt.Println("请输入正确的选项")
			return false
		}
	} else {
		fmt.Println("-------------------------------")
		fmt.Println("-    公聊      -      2       -")
		fmt.Println("-    私聊      -      3       -")
		fmt.Println("-   设定昵称   -      4       -")
		fmt.Println("- 查询在线用户 -      5       -")
		fmt.Println("-------------------------------")
		fmt.Print("请选择：")
		fmt.Scan(&flag)
		if flag >= 2 || flag < 6 {
			//fmt.Println("-----")
			this.flag = flag
			return true
		} else {
			fmt.Println("请输入正确的选项")
			return false
		}
	}
}

func (this *Cli) Login() {
	msg := "login"
	this.conn.Write([]byte(msg + "\n"))
	tmp := make([]byte, 4096)
	var id, pwd string
	for {
		fmt.Print("请输入你的账号：")
		fmt.Scanln(&id)
		fmt.Print("请输入你的密码：")
		fmt.Scanln(&pwd)
		this.conn.Write([]byte(id + "@" + pwd + "\n"))
		this.conn.Read(tmp)
		if string(tmp)[:5] == "true\n" {
			fmt.Println("登录成功")
			this.login = true
			//this.menu()
			break
		} else {
			fmt.Println("登录失败，请检查账号和密码后尝试")
		}
	}
}

func (this *Cli) Register() {
	var id, pwd string
	tmp := make([]byte, 4096)
	for {
		fmt.Print("请输入你要注册的id：")
		fmt.Scan(&id)
		fmt.Print("请输入密码：")
		fmt.Scan(&pwd)
		msg := "re|" + id + "|" + pwd + "\n"
		this.conn.Write([]byte(msg))
		this.conn.Read(tmp)
		fmt.Print(string(tmp))
		if string(tmp)[:4] == "err\n" {
			fmt.Println("此id已存在，请重新注册")
			continue
		} else if string(tmp)[:5] == "true\n" {
			fmt.Println("注册成功")
			break
		} else if string(tmp)[:6] == "false\n" {
			fmt.Println("注册失败，请检查后重新注册")
			continue
		}
	}
	this.Login()
}

func (this *Cli) SearchUser() {
	msg := "ls\n"
	if _, err := this.conn.Write([]byte(msg)); err != nil {
		fmt.Println("查询错误", err)
		return
	}
	var p string
	fmt.Println("如下是在线用户的昵称，输入“exit”退出")
	fmt.Scan(&p)
	for p != "exit" {
	}
}

func (this *Cli) SetName() {
	fmt.Println("请输入你要设定的新昵称：")
	var rname string
	fmt.Scan(&rname)
	this.Name = rname
	msg := "set|" + rname + "\n"
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("write err", err)
		return
	}
	fmt.Println("修改成功")
}

func (this *Cli) PublicChat() {
	fmt.Println("当前为公聊模式，请输入消息内容，输入“exit”退出")
	var msg string
	fmt.Scan(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			sMsg := msg + "\n"
			_, err := this.conn.Write([]byte(sMsg))
			if err != nil {
				fmt.Println("conn write err", err)
				return
			}
		}
		msg = ""
		//fmt.Println("当前为公聊模式，请输入消息内容，输入“exit”退出")
		fmt.Scan(&msg)
	}
}

func (this *Cli) Search() {
	msg := "ls\n"
	//fmt.Println("如下是你可以私聊的用户,请输入你要私聊的用户昵称：")
	if _, err := this.conn.Write([]byte(msg)); err != nil {
		fmt.Println("查询在线用户错误", err)
		return
	}
}

func (this *Cli) PraviteChat() {
	this.Search()
	//fmt.Print("请输入你要私聊的用户昵称：")
	var uname string

	//fmt.Println("当前在线用户如下所示，输入“exit”退出")
	fmt.Println("如下是你可以私聊的用户,输入“exit”退出,请输入你要私聊的用户昵称：")
	fmt.Scan(&uname)
	for uname != "exit" {
		fmt.Println("请输入你私聊的内容，输入“exit”退出")
		var msg string
		fmt.Scan(&msg)
		for msg != "exit" {
			if len(msg) != 0 {
				sendmsg := "to|" + uname + "|" + msg + "\n"
				_, err := this.conn.Write([]byte(sendmsg))
				if err != nil {
					fmt.Println("write err", err)
					break
				}
			}
			msg = ""
			fmt.Println("继续输入你想私聊的内容，输入“exit”退出")
			fmt.Scan(&msg)
		}
		uname = ""
		this.Search()
		fmt.Println("如下是你可以私聊的用户,输入“exit”退出,请输入你要私聊的用户昵称：")
		fmt.Scan(&uname)
	}
}

func (this *Cli) Run() {
	for this.flag != 0 {
		for this.menu() != true {
			//fmt.Println("======")
		}
		//fmt.Println("=====")
		switch this.flag {
		case 2:
			this.PublicChat()
		case 3:
			this.PraviteChat()
		case 4:
			this.SetName()
		case 5:
			this.SearchUser()
		}
	}
}

func main() {
	client := NewCli("192.168.153.128", 6000)
	//client.login = true
	//client.menu()
	//client.menu()
	//client.Login()
	for client.login != true {
		client.menu()
		if client.flag == 1 {
			client.Register()
			break
		}
		client.Login()
	}
	fmt.Println("------")
	go func() {
		io.Copy(os.Stdout, client.conn)
	}()
	client.Run()
	//client.Login()
}
