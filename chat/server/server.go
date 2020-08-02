package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/LeadNess/net-tools/chat/protocol"
)

type ChatServer interface {
	Listen(address string) error
	Broadcast(command interface{}) error
	ClientsUsernames() []string
	Start()
	Close() error
	Logs() chan string
	Clients() chan []*client
}

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
	logs     chan string
	clientsChan chan []*client
}

type client struct {
	Conn   net.Conn
	Name   string
	writer *protocol.CommandWriter
}

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		mutex: &sync.Mutex{},
		logs: make(chan string, 10),
		clientsChan: make(chan []*client),
	}
}

func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
		s.logs <- fmt.Sprintf("%s Listening on %v",
			time.Now().Format("15:04"), address)
	}
	return err
}

func (s *TcpChatServer) Close() error {
	return s.listener.Close()
}

func (s *TcpChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *TcpChatServer) accept(conn net.Conn) *client {
	s.logs <- fmt.Sprintf("%s Accepting connection from %v, total clients: %v",
		time.Now().Format("15:04"), conn.RemoteAddr().String(), len(s.clients)+1)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	client := &client{
		Name:   conn.RemoteAddr().String(),
		Conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}
	s.clients = append(s.clients, client)
	return client
}

func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	s.logs <- fmt.Sprintf("%s Closing connection from %v",
		time.Now().Format("15:04"), client.Conn.RemoteAddr().String())

	s.clientsChan <- s.clients
	go s.Broadcast(protocol.UsersCommand{
		Users: strings.Join(s.ClientsUsernames(), " "),
	})

	client.Conn.Close()
}

func (s *TcpChatServer) serve(client *client) {
	cmdReader := protocol.NewCommandReader(client.Conn)
	defer s.remove(client)
	for {
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF {
			s.logs <- fmt.Sprintf("%s Read error: %v",
				time.Now().Format("15:04"), err)
		}
		if cmd != nil {
			switch v := cmd.(type) {
				case protocol.SendCommand:
					go s.Broadcast(protocol.MessageCommand{
						Message: v.Message,
						Name:    client.Name,
					})
				case protocol.NameCommand:
					client.Name = v.Name
					s.clientsChan <- s.clients
					go s.Broadcast(protocol.UsersCommand{
						Users: strings.Join(s.ClientsUsernames(), " "),
					})
			}
		}
		if err == io.EOF {
			break
		}
	}
}

func (s *TcpChatServer) ClientsUsernames() []string {
	var users []string
	for _, client := range s.clients {
		users = append(users, client.Name)
	}
	return users
}

func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		if err := client.writer.Write(command); err != nil {
			log.Printf("Broadcast error: %v", err)
		}
	}
	return nil
}

func (s *TcpChatServer) Logs() chan string {
	return s.logs
}

func (s *TcpChatServer) Clients() chan []*client {
	return s.clientsChan
}