package main

func main() {
	server := NewServer("192.168.153.128", 6000)
	server.Start()
}
