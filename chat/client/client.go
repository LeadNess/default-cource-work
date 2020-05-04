package client

import (
	"default-cource-work/chat/protocol"
	"io"
	"log"
	"net"
	"strings"
)

type ChatClient interface {
	Dial(address string) error
	SendMessage(message string) error
	SetName(name string) error
	Start()
	Close()
	Incoming() chan protocol.MessageCommand
	ChatUsers() chan []string
}

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	incoming  chan protocol.MessageCommand
	users     chan []string
}

func NewClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
		users: make(chan []string),
	}
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
	}
	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)
	return err
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.cmdWriter.Write(protocol.SendCommand{
		Message: message,
	})
}

func (c *TcpChatClient) SetName(name string) error {
	return c.cmdWriter.Write(protocol.NameCommand{Name: name})
}

func (c * TcpChatClient) Incoming() chan protocol.MessageCommand  {
	return c.incoming
}

func (c * TcpChatClient) ChatUsers() chan []string  {
	return c.users
}

func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
		}
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			case protocol.UsersCommand:
				c.users <- strings.Split(v.Users, " ")
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}