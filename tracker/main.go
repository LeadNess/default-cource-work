package main

import (
	"bufio"
	"bytes"
	"default-cource-work/tracker/tracker"
	"fmt"
	"github.com/google/gopacket"
	"log"
	"os"
)

func main()  {
	device, err := tracker.GetDeviceName()
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	var filter string
	fmt.Print("Capture filter: ")
	filter, _ = reader.ReadString('\n')

	var maxLen int16
	fmt.Print("TCP payload data max len: ")
	buf, _ := reader.ReadString('\n')
	if maxLen, err = tracker.ParseInt16(buf); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	var snapLen int32 = 1600
	handler, err := tracker.GetHandle(device, filter, snapLen)
	if err != nil {
		log.Printf("Error on opening device: %v", err)
		os.Exit(1)
	}
	modifiedTrafficWriter, err := tracker.GetPcapFileWriter(uint32(snapLen), "Pcap filename for modified traffic: ")
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	allowedTrafficWriter, err := tracker.GetPcapFileWriter(uint32(snapLen), "Pcap filename for allowed traffic: ")
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	cfgFileName := "tracker.cfg"
	var cfgMap map[string]interface{}
	if cfgMap, err = tracker.GetConfig(cfgFileName); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	var A, B, D *tracker.Computer
	if A, err = tracker.NewComputer(cfgMap["ipA"].(string), cfgMap["portA"].(uint16)); err != nil {
		log.Printf("Error on creating A computer instance: %v", err)
		os.Exit(1)
	}
	if B, err = tracker.NewComputer(cfgMap["ipC"].(string), cfgMap["portC"].(uint16)); err != nil {
		log.Printf("Error on creating B computer instance: %v", err)
		os.Exit(1)
	}
	if D, err = tracker.NewComputer(cfgMap["ipD"].(string), cfgMap["portD"].(uint16)); err != nil {
		log.Printf("Error on creating D computer instance: %v", err)
		os.Exit(1)
	}

	for packet := range gopacket.NewPacketSource(handler, handler.LinkType()).Packets() {
		if packet.TransportLayer() != nil && packet.TransportLayer().LayerType().String() == "UDP" {
			p := tracker.NewPacket(packet.Data())
			if bytes.Compare(p.SrcIPv4, A.IPv4) == 0 && bytes.Compare(p.SrcPort, A.Port) == 0 {
				p.SrcIPv4 = B.IPv4
				p.SrcPort = B.Port
				if err = handler.WritePacketData(p.Data); err != nil {
					log.Printf("Error on sending packet: %v", err)
				} else {
					if err = modifiedTrafficWriter.WritePacket(packet.Metadata().CaptureInfo, p.Data); err != nil {
						log.Printf("Error on writing packet to pcap file: %v", err)
					}
				}
			}
			if bytes.Compare(p.SrcIPv4, D.IPv4) == 0 && bytes.Compare(p.SrcPort, D.Port) == 0 {
				p.DstIPv4 = A.IPv4
				p.DstPort = A.Port
				if err = handler.WritePacketData(p.Data); err != nil {
					log.Printf("Error on sending packet: %v", err)
				} else {
					if err = modifiedTrafficWriter.WritePacket(packet.Metadata().CaptureInfo, p.Data); err != nil {
						log.Printf("Error on writing packet to pcap file: %v", err)
					}
				}
			}
		}
		if packet.TransportLayer() != nil && packet.TransportLayer().LayerType().String() == "TCP" {
			if len(packet.TransportLayer().LayerPayload()) < int(maxLen) {
				if err = allowedTrafficWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
					log.Printf("Error on writing packet to pcap file: %v", err)
				}
			}
		}
	}
}