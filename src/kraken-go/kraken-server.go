package main

import (
	"os"
	"fmt"
	"./tcp"
	"./errors"
)

func main() {
	sock := tcp.CreateSocket()

	sock.OnMessage(func(client *tcp.SocketClient, message *tcp.SocketMessage) {
		fmt.Printf("%s\n", message.GetRecord().Get(tcp.MESSAGE_TEXT))
	})

	err := sock.Listen("127.0.0.1", "9080")
	errors.Log(err)

	sock.Lock()
	sock.Close()

	os.Exit(0)
}
