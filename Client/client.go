package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct{
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string,serverPort int)*Client{
	//创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 9999,
	}

	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil{
		fmt.Println("net.Dial err",err)
		return nil
	}

	client.conn = conn

	return client
}

func (client *Client) menu() bool{
	var flag int
	
	fmt.Println("1->公聊模式")
	fmt.Println("2->私聊模式")
	fmt.Println("3->更改用户名")
	fmt.Println("0->退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3{
		client.flag = flag
		return true
	}else{
		fmt.Println("请输入正确的选项")
		return false
	}
}

var serverIp string
var serverPort int

func (client *Client) UpdateName() bool{
	fmt.Println("请输入你想修改后的用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|"+client.Name+"\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn.Write err",err)
		return false
	}

	return true
}

func (client *Client) SelectUsers(){
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("select err",err)
		return
	}
}


func (client *Client) PrivateChat(){
	var sMsg string
	var sName string
	client.SelectUsers()
	fmt.Println("请输入聊天对象(用户名),exit退出")
	fmt.Scanln(&sName)

	for sName != "exit"{
		fmt.Println("请输入聊天内容，输入exit退出私聊")
		fmt.Scanln(&sMsg)
		for sMsg != "exit"{
                
	                if len(sMsg) != 0{
        	                sendMsg := "to|" + sName + "|" + sMsg + "\n"
                	        _, err := client.conn.Write([]byte(sendMsg))
                        	if err != nil{
                        	        fmt.Println("conn Write err",err)
                        	        break
                       		 }
               		 }

             		sMsg = ""
                	fmt.Println("请输入聊天内容，输入exit退出私聊")
                	fmt.Scanln(&sMsg)
		}
		client.SelectUsers()
        fmt.Println("请输入聊天对象(用户名),exit退出")
        fmt.Scanln(&sName)
	}
	
}

func (client *Client) PublicChat(){
	var sMsg string
	fmt.Println("请输入聊天内容，输入exit退出")
	fmt.Scanln(&sMsg)

	for sMsg != "exit"{
		
		if len(sMsg) != 0{
			sendMsg := sMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil{
				fmt.Println("conn Write err",err)
				break
			}
		}

		sMsg = ""
		fmt.Println("请输入聊天内容，输入exit退出")
		fmt.Scanln(&sMsg)
	}
}

func (client *Client) DealServerMsg(){
	io.Copy(os.Stdout, client.conn)
	/*相当于
	for{
		buf := make()
		client.conn.Read(buf)
		fmt.Println(buff)
	}*/
	
}

func (client *Client) Run(){
	for client.flag != 0{
		for client.menu() != true{
		}
		
		switch client.flag{
		case 1:
			client.PublicChat()
		//
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		}
	}
}

func Login()

func init(){
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器的端口(默认是8888)")
}

func main(){
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil{
		fmt.Println("连接服务器失败")
		return
	}

	go client.DealServerMsg()

	client.Run()
}
