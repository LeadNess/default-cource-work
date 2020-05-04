package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type SendCommand struct {
	Message string
}

type NameCommand struct {
	Name string

}

type MessageCommand struct {
	Name    string
	Message string
}

type UsersCommand struct {
	Users string
}

type UnknownCommand interface {
	Error() string
}

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	var err error
	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v\n", v.Name, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	case UsersCommand:
		err = w.writeString(fmt.Sprintf("USERS %v\n", v.Users))
	}
	return err
}

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (interface{}, error) {
	bufstr, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	bufslice := strings.Split(bufstr[:len(bufstr)-1], " ")
	commandName := bufslice[0]
	switch commandName {
	case "SEND":
		message := strings.Join(bufslice[1:], " ")
		return SendCommand{
			message,
		}, nil
	case "MESSAGE":
		user := bufslice[1]
		message := strings.Join(bufslice[2:], " ")
		return MessageCommand{
			user,
			message,
		}, nil
	case "NAME":
		name := bufslice[1]
		return NameCommand{
			name,
		}, nil
	case "USERS":
		users := strings.Join(bufslice[1:], " ")
		return UsersCommand{
			users,
		}, nil
	}
	log.Printf("Unknown command: %v", commandName)
	return nil, errors.New("unknown command")
}
