package tracker

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

type Computer struct {
	Port []byte
	IPv4 []byte
}

func NewComputer(ipv4 string, port uint16) (*Computer, error) {
	buf := strings.Split(ipv4, ".")
	if len(buf) != 4 {
		return nil, errors.New("incorrect ipv4 string format")
	}
	ipBytes := []byte{ParseByte(buf[0]), ParseByte(buf[1]), ParseByte(buf[2]), ParseByte(buf[3])}
	portBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(portBytes, port)
	return &Computer{
		Port: portBytes,
		IPv4: ipBytes,
	}, nil
}

type PacketData struct {
	Data    []byte
	DstMAC  []byte
	SrcMAC  []byte
	DstIPv4 []byte
	SrcIPv4 []byte
	DstPort []byte
	SrcPort []byte
}

func NewPacket(pkData []byte) *PacketData {
	return &PacketData{
		Data:    pkData,
		DstMAC:  pkData[0:6],
		SrcMAC:  pkData[6:12],
		DstIPv4: pkData[30:34],
		SrcIPv4: pkData[26:30],
		DstPort: pkData[36:38],
		SrcPort: pkData[34:36],
	}
}

func (p *PacketData) SetSrcIPv4(ip []byte)  {
	p.SrcIPv4 =  ip
	p.Data[26] = ip[0]
	p.Data[27] = ip[1]
	p.Data[28] = ip[2]
	p.Data[29] = ip[3]
}

func (p *PacketData) SetDstIPv4(ip []byte)  {
	p.DstIPv4 =  ip
	p.Data[30] = ip[0]
	p.Data[31] = ip[1]
	p.Data[32] = ip[2]
	p.Data[33] = ip[3]
}

func (p *PacketData) SetSrcPort(port []byte)  {
	p.SrcPort =  port
	p.Data[34] = port[0]
	p.Data[35] = port[1]
}

func (p *PacketData) SetDstPort(port []byte)  {
	p.DstPort =  port
	p.Data[36] = port[0]
	p.Data[37] = port[1]
}

func GetHandle(device, filter string, snapLen int32) (*pcap.Handle, error) {
	if handle, err := pcap.OpenLive(device, snapLen, true, pcap.BlockForever); err != nil {
		return nil, err
	} else if err := handle.SetBPFFilter(filter); err != nil {
		return nil, err
	} else {
		return handle, nil
	}
}

func ParseByte(s string) byte {
	i, _ := strconv.ParseInt(s, 10, 16)
	b := byte(i)
	return b
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

func ParseUint16(s string) (uint16, error) {
	var buf strings.Builder
	str := []rune(s)
	for c := range str {
		if unicode.IsNumber(str[c]) {
			buf.WriteRune(str[c])
		}
	}
	num, err := strconv.ParseUint(buf.String(), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(num), nil
}

func GetDeviceName() (string, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}
	fmt.Println("+---+---------------------+-----------------+----------+")
	fmt.Printf("|   | %-20s| %-16s| %-9s|\n", "Name", "IP", "Netmask")
	fmt.Println("+---+---------------------+-----------------+----------+")
	numToDev := make(map[int16]string, len(devices))
	var counter int16
	for _, dev := range devices {
		if len(dev.Addresses) > 0 {
			counter++
			fmt.Printf("| %-2d| %-20s| %-16s| %-9s|\n",
				counter, dev.Name, dev.Addresses[0].IP, dev.Addresses[0].Netmask)
			fmt.Println("+---+---------------------+-----------------+----------+")
			numToDev[counter] = dev.Name
		}
	}
	var devNum int16
	fmt.Print("Device: ")
	buf, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if devNum, err = ParseInt16(buf); err != nil {
		return "", err
	}
	return numToDev[devNum], nil
}

func GetPcapFileWriter(snapLen uint32, promptStr string) (*pcapgo.Writer, error) {
	var filename string
	fmt.Print(promptStr)
	filename, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	if runtime.GOOS == "linux" {
		filename = filename[:len(filename)-1]
	} else { //windows
		filename = filename[:len(filename)-2]
	}
	f, _ := os.Create(fmt.Sprintf("%s.pcap", filename))
	w := pcapgo.NewWriter(f)
	if err := w.WriteFileHeader(snapLen, layers.LinkTypeEthernet); err != nil {
		return nil, err
	}
	return w, nil
}

func GetConfig(cfgFileName string) (map[string]interface{}, error) {
	cfgMap := map[string]interface{}{
		"ipA": "",
		"portA": 0,
		"ipB": "",
		"portB": 0,
		"ipD": "",
		"portD": 0,
	}
	cfg, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(cfg), "\n")

	cfgMap["ipA"] = strings.Split(lines[0], " = ")[1]
	cfgMap["portA"], err = ParseUint16(strings.Split(lines[1], " = ")[1])

	cfgMap["ipB"] = strings.Split(lines[2], " = ")[1]
	cfgMap["portB"], err = ParseUint16(strings.Split(lines[3], " = ")[1])

	cfgMap["ipD"] = strings.Split(lines[4], " = ")[1]
	cfgMap["portD"], err = ParseUint16(strings.Split(lines[5], " = ")[1])

	return cfgMap, err
}