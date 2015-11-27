package tcp

import (
	"net"
	"bufio"
	"strings"
	"../storage"
	"../errors"
)

const (
	SOCKET_SUCCESS				 int = 1
	SOCKET_CLOSED_COUT			 int = 2
	SOCKET_CLOSED_CIN			 int = 3
	SOCKET_ERR_NOT_STARTED		 int = 4
	SOCKET_ERR_NOT_CONNECTED	 int = 5
	SOCKET_ERR_ALREADY_CONNECTED int = 6
)

const (
	SOCKET_MESSAGE			 	 string = "MSG"
	SOCKET_COMMAND			 	 string = "CMD"
)

const (
	COMMAND_EXIT				 string = "EXIT"
)

const (
	MESSAGE_TEXT				 string = "TXT"
)

//--------------------------------------------------------------------------------------------------------------------//
/**
 * SocketMessage class
 */
type SocketMessage struct {
	Cmd		string
	Val		*storage.DataRecord
}

/**
 * SocketMessage constructor
 */
func CreateSocketMessage(cmd string, val *storage.DataRecord) *SocketMessage {
	return &SocketMessage{cmd, val}
}

/**
 * SocketMessage.GetCmd() string
 */
func (message *SocketMessage) GetCmd() string {
	return message.Cmd
}

/**
 * SocketMessage.GetRecord() *storage.DataRecord
 */
func (message *SocketMessage) GetRecord() *storage.DataRecord {
	return message.Val
}
/**
 * SocketMessage.ToString() string
 */
func (message *SocketMessage) ToString() string {
	return "[" + message.Cmd + "]" + message.Val.ToString()
}

/**
 * SocketMessage.FromString(string) *SocketMessage
 */
func (message *SocketMessage) FromString(line string) *SocketMessage {
	parts := strings.SplitN(line, "]", 2)
	parts[0] = parts[0][1:]

	message.Cmd = parts[0]
	message.Val = storage.CreateDataRecord().FromString(parts[1])

	return message
}

//--------------------------------------------------------------------------------------------------------------------//
/**
 * SocketEvents
 */
type SocketEvents struct {
	Start       		func()
	Stop        		func()
	Message     		func(c *SocketClient, s *SocketMessage)
	ClientStart			func(c *SocketClient)
	ClientStop			func(c *SocketClient)
}

/**
 * SocketEvents constructor
 */
func CreateSocketEvents() *SocketEvents {
	events := &SocketEvents{}

	events.Start    	= func() {}
	events.Stop     	= func() {}
	events.Message  	= func(c *SocketClient, s *SocketMessage) {}
	events.ClientStart	= func(c *SocketClient) {}
	events.ClientStop	= func(c *SocketClient) {}

	return events
}

//--------------------------------------------------------------------------------------------------------------------//
/**
 * SocketClientFlags class
 */
type SocketClientFlags struct {
	IsListening		bool
}

/**
 * SocketClientFlags constructor
 */
func CreateSocketClientFlags() *SocketClientFlags {
	flags := &SocketClientFlags{}

	flags.IsListening = true

	return flags
}

//--------------------------------------------------------------------------------------------------------------------//
/*
 * SocketClient
 */
type SocketClient struct {
	sock		*Socket
	conn		net.Conn
	cin			*bufio.Writer
	cout		*bufio.Reader
	flags		*SocketClientFlags
}

/**
 * SocketClient constructor
 */
func CreateSocketClient(sock *Socket, conn net.Conn) *SocketClient {
	client := &SocketClient{}

	client.sock	  = sock
	client.conn   = conn
	client.cin	  = bufio.NewWriter(conn)
	client.cout	  = bufio.NewReader(conn)
	client.flags  = CreateSocketClientFlags()

	return client
}

/**
 * SocketClient.Close()
 */
func (c *SocketClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

/**
 * SocketClient.Reset()
 */
func (c *SocketClient) Reset() {
	c.flags = CreateSocketClientFlags()
}

/**
 * SocketClient.ReadMessage() (*SocketMessage, errors.Error)
 */
func (c *SocketClient) ReadMessage() (*SocketMessage, errors.Error) {
	line, err := c.cout.ReadString('\n')
	if err != nil {
		return nil, errors.New(SOCKET_CLOSED_COUT, err.Error())
	}

	msg := strings.TrimSpace(line)

	return CreateSocketMessage("",nil).FromString(msg), nil
}

/**
 * SocketClient.WriteMessage(*SocketMessage) errors.Error
 */
func (c *SocketClient) WriteMessage(m *SocketMessage) errors.Error {
	_, err := c.cin.WriteString(m.ToString() + "\n")
	if err != nil {
		return errors.New(SOCKET_CLOSED_CIN, err.Error())
	}

	c.cin.Flush()

	return nil
}

//--------------------------------------------------------------------------------------------------------------------//
/*
 * SocketConn class
 */
type SocketConn struct {
	Client		*SocketClient
	Clients		[]*SocketClient
	Accept		func() (net.Conn, error)
	Close		func() (error)
}

//--------------------------------------------------------------------------------------------------------------------//
/*
 * Socket
 */
type Socket struct {
	IsConnected bool
	Type        string
	Host        string
	Port        string
	Sync        chan *SocketMessage
	Cin         chan *SocketMessage
	Cout        chan *SocketMessage
	Conn        *SocketConn
	Events      *SocketEvents
}

/**
 * Socket constructor
 */
func CreateSocket() *Socket {
	sock := &Socket{}

	sock.IsConnected  = false
	sock.Type         = "tcp"
	sock.Host         = ""
	sock.Port         = ""
	sock.Sync         = make(chan *SocketMessage)
	sock.Cin          = make(chan *SocketMessage)
	sock.Cout         = make(chan *SocketMessage)
	sock.Conn		  = nil
	sock.Events       = CreateSocketEvents()

	return sock
}

/**
 * Socket.OnStart(func())
 */
func (sock *Socket) OnStart(f func()) {
	sock.Events.Start = f
}

/**
 * Socket.OnStop(func())
 */
func (sock *Socket) OnStop(f func()) {
	sock.Events.Stop = f
}

/**
 * Socket.OnMessage(func(*SocketClient, *SocketMessage))
 */
func (sock *Socket) OnMessage(f func(c *SocketClient, s *SocketMessage)) {
	sock.Events.Message = f
}

/**
 * Socket.OnClientStart(func(*SocketClient))
 */
func (sock *Socket) OnClientStart(f func(c *SocketClient)) {
	sock.Events.ClientStart = f
}

/**
 * Socket.OnClientStop(func(*SocketClient))
 */
func (sock *Socket) OnClientStop(f func(c *SocketClient)) {
	sock.Events.ClientStop = f
}

/**
 * Socket.Listen(string, string) errors.Error
 */
func (sock *Socket) Listen(host string, port string) errors.Error {
	if sock.IsConnected {
		return errors.New(SOCKET_ERR_ALREADY_CONNECTED, "Socket connection has been already estabilished.")
	}

	sock.Host = host
	sock.Port = port

	conn, err := net.Listen(sock.Type, sock.Host + ":" + sock.Port)
	if err != nil {
		return errors.New(SOCKET_ERR_NOT_STARTED, err.Error())
	}

	sock.Conn = &SocketConn{}
	sock.Conn.Client = nil
	sock.Conn.Accept = conn.Accept
	sock.Conn.Close  = conn.Close

	sock.IsConnected = true
	sock.Events.Start()

	sock.Accept()

	return nil
}

/**
 * Socket.Connect(string, string) errors.Error
 */
func (sock *Socket) Connect(host string, port string) errors.Error {
	if sock.IsConnected {
		return errors.New(SOCKET_ERR_ALREADY_CONNECTED, "Socket connection has been already estabilished.")
	}

	sock.Host = host
	sock.Port = port

	conn, err := net.Dial(sock.Type, sock.Host + ":" + sock.Port)
	if err != nil {
		return errors.New(SOCKET_ERR_NOT_STARTED, err.Error())
	}

	sock.Conn = &SocketConn{}
	sock.Conn.Client = CreateSocketClient(sock, conn)
	sock.Conn.Accept = func() (c net.Conn, err error) {
		return nil, nil
	}
	sock.Conn.Close  = conn.Close

	sock.IsConnected = true
	sock.Events.Start()

	sock.Accept()

	return nil
}

/**
 * Socket.Close() errors.Error
 */
func (sock *Socket) Close() errors.Error {
	if !sock.IsConnected {
		return errors.New(SOCKET_ERR_NOT_CONNECTED, "Cannot close null connection.")
	}

	sock.Conn.Close()
	sock.IsConnected = false
	sock.Events.Stop()

	return nil
}

/**
 * Socket.Accept() errors.Error
 */
func (sock *Socket) Accept() errors.Error {
	if !sock.IsConnected {
		return errors.New(SOCKET_ERR_NOT_CONNECTED, "Cannot accept connection from closed connection.")
	}

	// listener
	if sock.Conn.Client == nil {
		go func() {
			for sock.IsConnected {
				sock.acceptConnection()
			}
		}()

	// connector
	} else {
		go func() {
			if sock.IsConnected {
				sock.readConnection(sock.Conn.Client)
			}
		}()
	}

	return nil
}

/**
 * Socket.Lock() errors.Error
 */
func (sock *Socket) Lock() errors.Error {
	if !sock.IsConnected {
		return errors.New(SOCKET_ERR_NOT_CONNECTED, "Cannot lock closed connection.")
	}

	sock.freeze()
	sock.Close()

	return nil
}

/**
 * Socket.Unlock() errors.Error
 */
func (sock *Socket) Unlock() errors.Error {
	record := storage.CreateDataRecord()
	record.Set("Command", COMMAND_EXIT)
	sock.Sync <- CreateSocketMessage(SOCKET_COMMAND, record)

	return nil
}

/**
 * Socket.WriteMessage(*SocketClient, *SocketMessage) errors.Error
 */
func (sock *Socket) WriteMessage(c *SocketClient, m *SocketMessage) errors.Error {
	if sock.IsConnected && c != nil {
		return c.WriteMessage(m)
	}
	return nil
}

/**
 * Socket.ReadMessage(*SocketClient) (*SocketMessage, errors.Errror)
 */
func (sock *Socket) ReadMessage(c *SocketClient) (*SocketMessage, errors.Error) {
	return c.ReadMessage()
}

/**
 * Socket.acceptConnection()
 */
func (sock *Socket) acceptConnection() {
	listener, _ := sock.Conn.Accept()

	client := CreateSocketClient(sock, listener)
	sock.Conn.Clients = append(sock.Conn.Clients, client)
	sock.Events.ClientStart(client)

	go func(sock *Socket, client *SocketClient) {
		defer sock.closeConnection(client)

		record := storage.CreateDataRecord()
		record.Set(MESSAGE_TEXT, "Response")

		sock.WriteMessage(client, CreateSocketMessage(SOCKET_MESSAGE, record))

		sock.readConnection(client)
	}(sock, client)
}

/**
 * Socket.closeConnection(*SocketClient)
 */
func (sock *Socket) closeConnection(c *SocketClient) {
	sock.Events.ClientStop(c)
	c.Close()
}

/**
 * Socket.readConnection(*SocketClient)
 */
func (sock *Socket) readConnection(c *SocketClient) {
	stopFlag := false
	for !stopFlag && sock.IsConnected {
		message, _ := sock.ReadMessage(c)

		if message == nil {
			stopFlag = true
		} else {
			sock.Events.Message(c, message)
		}
	}
}

/**
 * Socket.freeze()
 */
func (sock *Socket) freeze() {
	stopFlag := false

	for !stopFlag && sock.IsConnected {
		message := <-sock.Sync

		cmd 	:= message.GetCmd()
		record  := message.GetRecord()

		if cmd == SOCKET_COMMAND && record.Exists("Command") {
			switch record.Get("Command")  {
				case COMMAND_EXIT:
					stopFlag = true
			}
		}
	}
}
