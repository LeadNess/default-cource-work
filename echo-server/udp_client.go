package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func parseArgsServer() (address string, err error) {
	const cfgFilename = "upd_client.cfg"
	var port int
	if len(os.Args) == 1 {
		cfg, err := ioutil.ReadFile(cfgFilename)
		if err != nil {
			return address, err
		}
		if _, err := fmt.Sscanf(string(cfg), "address = %s\nport = %d", &address, &port); err != nil {
			return address, err
		}
	}
	if len(os.Args) == 3 {
		address = os.Args[1]
		if port, err = strconv.Atoi(os.Args[2]); err != nil {
			return address, err
		}
	}
	if len(os.Args) == 2 || len(os.Args) > 3 {
		log.Fatalf("Usage %s: <address> <port>\n", filepath.Base(os.Args[0]))
	}

	address = fmt.Sprintf("%s:%d", address, port)
	return address, nil
}

func main() {
	address, err := parseArgsServer()
	if err != nil {
		log.Fatal(err)
	}

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	const maxBufferSize = 1024
	buffer := make([]byte, maxBufferSize)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Connected to server %s...\n", address)

	for {
		fmt.Print("Send message: ")
		message, _ := reader.ReadString('\n')
		if _, err := fmt.Fprint(conn, message); err != nil {
			log.Fatal(err)
		}

		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received message: %s", buffer[:n])
	}
}