package protocol

import (
	"fmt"
	"io"
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

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

type CommandReader struct {
	reader io.Reader
}

func (r CommandReader) Read() (p []byte, err error) {
	buffer := make([]byte, 1024)
	n, err := r.reader.Read(buffer)
	if err != nil {
		return p, err
	}
	return buffer[:n], nil
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: reader,
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
	}
	return err
}

