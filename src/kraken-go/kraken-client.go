package main

import (
	"os"
	"fmt"
	"time"
	"./tcp"
	"./storage"
	"./errors"
)

func main() {
	sock := tcp.CreateSocket()

	sock.OnMessage(func(client *tcp.SocketClient, message *tcp.SocketMessage) {
		fmt.Printf("%s\n", message.GetRecord().Get(tcp.MESSAGE_TEXT))
	})

	err := sock.Connect("127.0.0.1", "9080")
	errors.Log(err)

	lock := make(chan bool)

	go func(sock *tcp.Socket, lock chan bool) {
		stopFlag := false
		for !stopFlag {
			record := storage.CreateDataRecord()
			record.Set(tcp.MESSAGE_TEXT, "Request")

			err := sock.WriteMessage(sock.Conn.Client, tcp.CreateSocketMessage(tcp.SOCKET_MESSAGE, record))
			if err != nil {
				sock.Close()
				stopFlag = true
				lock <- true

			} else {
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}(sock, lock)

	exitcnt := 1
	for {
		<-lock
		exitcnt = exitcnt - 1
		if exitcnt == 0 {
			break
		}
	}

//	sock.Lock()
//	sock.Close()

	os.Exit(0)
}
