package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func parseArgsClient() (address string, port int, err error) {
	const cfgFilename = "upd_server.cfg"
	if len(os.Args) == 1 {
		cfg, err := ioutil.ReadFile(cfgFilename)
		if err != nil {
			return address, port, err
		}
		if _, err := fmt.Sscanf(string(cfg), "port = %d", &port); err != nil {
			return address, port, err
		}
	}
	if len(os.Args) == 2 {
		if port, err = strconv.Atoi(os.Args[1]); err != nil {
			return address, port, err
		}
	}
	if len(os.Args) > 2 {
		log.Fatalf("Usage %s: <port>\n", filepath.Base(os.Args[0]))
	}

	address = fmt.Sprintf("%s:%d", address, port)
	return address, port, nil
}

func main() {
	address, port, err := parseArgsClient()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Launching server on port %d...\n", port)

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	const maxBufferSize = 1024
	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := pc.ReadFrom(buffer)
		if n == 0 {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Received from %s: %s",
			addr.String(), buffer[:n])

		n, err = pc.WriteTo(buffer[:n], addr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send to %s: %s",
			addr.String(), buffer[:n])
	}
}