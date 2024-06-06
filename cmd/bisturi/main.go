package main

import (
	"flag"
	"log"
	"net"

	"github.com/NamelessOne91/bisturi/socket"
)

var iface = flag.String("i", "eth0", "The network interface to listen to")
var protocol = flag.String("p", "all", "Consider only packets for this protocol")

func main() {
	flag.Parse()

	// retrieve the network interface
	networkInterface, err := net.InterfaceByName(*iface)
	if err != nil {
		log.Fatalf("Failed to get interface by name: %v", err)
	}

	// SYS_SOCKET syscall
	rs, err := socket.NewRawSocket(*protocol)
	if err != nil {
		log.Fatalf("Failed to open raw socket: %v", err)
	}
	defer rs.Close()

	// bind the socket to the network interface
	rs.Bind(*networkInterface)
	if err != nil {
		log.Fatalf("Failed to bind socket: %v", err)
	}
	log.Printf("listening for %s packets on interface: %s\n", *protocol, networkInterface.Name)

	// SYS_RECVFROM syscall
	rs.ReadPackets()
}
