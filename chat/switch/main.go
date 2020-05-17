package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Cell struct {
	Port int16
	AddTime time.Time
}

type SwitchingTable map[string]Cell

func NewSwitchingTable(size uint16) SwitchingTable {
	return make(map[string]Cell, size)
}

func (table SwitchingTable) ContainsMAC(mac []byte) bool {
	 key := string(mac)
	 for cellKey := range table {
	 	if cellKey == key {
	 		return true
		}
	 }
	 return false
}

type EthHeader struct {
	DstMAC  []byte
	SrcMAC  []byte
	Type []byte
}

func ParseEthHeader(bytes []byte) (*EthHeader, error) {
	if len(bytes) != 14 {
		return nil, errors.New("incorrect link layer bytes")
	}
	return &EthHeader{
		DstMAC: bytes[0:6],
		SrcMAC: bytes[6:12],
		Type: bytes[12:],
	}, nil
}

func ParseInt16(s string) (int16, error) {
	var buf strings.Builder
	str := []rune(s)
	for c := range str {
		if unicode.IsNumber(str[c]) {
			buf.WriteRune(str[c])
		}
	}
	num, err := strconv.ParseInt(buf.String(), 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(num), nil
}

func ParseFloat(s string) (float64, error) {
	var buf strings.Builder
	str := []rune(s)
	for c := range str {
		if unicode.IsNumber(str[c]) || str[c] == ',' || str[c] == '.' {
			buf.WriteRune(str[c])
		}
	}
	return strconv.ParseFloat(buf.String(), 16)
}

func GetHandle(device, filter string) (*pcap.Handle, error) {
	if handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever); err != nil {
		return nil, err
	} else if err := handle.SetBPFFilter(filter); err != nil {
		return nil, err
	} else {
		return handle, nil
	}
}

func PrintSwitchingTable(table SwitchingTable)  {
	fmt.Println("+------+-------------------+")
	fmt.Println("| Port | MAC               |")
	fmt.Println("+------+-------------------+")
	for mac, cell := range table {
		fmt.Printf("| %d    | %02X:%02X:%02X:%02X:%02X:%02X |\n",
			cell.Port, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
		fmt.Println("+------+-------------------+")
	}
}

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Printf("Error on finding all devices: %v", err)
		os.Exit(1)
	}
	fmt.Println("+---+---------------------+-----------------+----------+")
	fmt.Printf("|   | %-20s| %-16s| %-9s|\n", "Name", "IP", "Netmask")
	fmt.Println("+---+---------------------+-----------------+----------+")
	numToDev := make(map[int16]string, len(devices))
	var counter int16
	for _, dev := range devices {
		if len(dev.Addresses) > 0 {
			counter++
			fmt.Printf("| %-2d| %-20s| %-16s| %-9s|\n", counter, dev.Name, dev.Addresses[0].IP, dev.Addresses[0].Netmask)
			fmt.Println("+---+---------------------+-----------------+----------+")
			numToDev[counter] = dev.Name
		}
	}
	var firstDevNum, secondDevNum int16
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Device 1: ")
	buf, _ := reader.ReadString('\n')
	if firstDevNum, err = ParseInt16(buf); err != nil {
		log.Printf("Error on parsing string to int: %v", err)
		os.Exit(1)
	}
	fmt.Print("Device 2: ")
	buf, _ = reader.ReadString('\n')
	if secondDevNum, err = ParseInt16(buf); err != nil {
		log.Printf("Error on parsing string to int: %v", err)
		os.Exit(1)
	}

	var cellTTL float64
	fmt.Print("Cells TTL: ")
	buf, _ = reader.ReadString('\n')
	if cellTTL, err = ParseFloat(buf); err != nil {
		log.Printf("Error on parsing string to float: %v", err)
		os.Exit(1)
	}

	var filter string
	fmt.Print("Capture filter: ")
	filter, _ = reader.ReadString('\n')

	firstPortHandler, err := GetHandle(numToDev[firstDevNum], filter)
	if err != nil {
		log.Printf("Error on opening first device: %v", err)
		os.Exit(1)
	}
	secondPortHandler, err := GetHandle(numToDev[secondDevNum], filter)
	if err != nil {
		log.Printf("Error on opening second device: %v", err)
		os.Exit(1)
	}

	swTable := NewSwitchingTable(2)
	var firstPortCounter, secondPortCounter int64

	go func() {
		for packet := range gopacket.NewPacketSource(firstPortHandler, firstPortHandler.LinkType()).Packets() {
			eth, err := ParseEthHeader(packet.LinkLayer().LayerContents())
			if err != nil {
				log.Printf("Error on parsing packet: %v", err)
				continue
			}
			if swTable.ContainsMAC(eth.SrcMAC) {
				if swTable[string(eth.SrcMAC)].Port == 1 {
					//fmt.Printf("IF PORT == 1 %d  -  %v\n\n", 1, packet)
					if err = secondPortHandler.WritePacketData(packet.Data()); err != nil {
						log.Printf("Error on sending packet: %v", err)
						continue
					}
				}
				firstPortCounter++
			} else {
				swTable[string(eth.SrcMAC)] = Cell{1, time.Now()}
				if err = secondPortHandler.WritePacketData(packet.Data()); err != nil {
					log.Printf("Error on sending packet: %v", err)
					continue
				}
				firstPortCounter++
			}
		}
	}()

	go func() {
		for packet := range gopacket.NewPacketSource(secondPortHandler, secondPortHandler.LinkType()).Packets() {
			eth, err := ParseEthHeader(packet.LinkLayer().LayerContents())
			if err != nil {
				log.Printf("Error on parsing packet: %v", err)
				continue
			}
			if swTable.ContainsMAC(eth.SrcMAC) {
				if swTable[string(eth.SrcMAC)].Port == 2 {
					if err = firstPortHandler.WritePacketData(packet.Data()); err != nil {
						log.Printf("Error on sending packet: %v", err)
						continue
					}
					secondPortCounter++
				}
			} else {
				swTable[string(eth.SrcMAC)] = Cell{2, time.Now()}
				if err = firstPortHandler.WritePacketData(packet.Data()); err != nil {
					log.Printf("Error on sending packet: %v", err)
					continue
				}
				secondPortCounter++
			}
		}
	}()

	go func() {
		for {
			for mac, cell := range swTable {
				now := time.Now()
				if now.Sub(cell.AddTime).Seconds() > cellTTL {
					delete(swTable, mac)
				}
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()

	for {
		fmt.Println("\n1. Print switching table")
		fmt.Println("2. Print packages count")
		fmt.Println("3. Exit")
		fmt.Print("Option: ")
		var option int16
		buf, _ := reader.ReadString('\n')
		if option, err = ParseInt16(buf); err != nil {
			log.Printf("Error on parsing option: %v", err)
			continue
		}
		switch option {
		case 1:
			PrintSwitchingTable(swTable)
		case 2:
			fmt.Printf("\nFirst port: %d\n", firstPortCounter)
			fmt.Printf("Second port: %d\n", secondPortCounter)
		case 3:
			fmt.Println("Total count")
			fmt.Printf("\nFirst port: %d\n", firstPortCounter)
			fmt.Printf("Second port: %d\n", secondPortCounter)
			fmt.Println("Bye!")
			return
		}
	}
}