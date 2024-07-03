package main

import (
	"log"

	tui "github.com/NamelessOne91/bisturi/tui/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.NewStartMenuModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}

	// // SYS_SOCKET syscall
	// rs, err := sockets.NewRawSocket(protocol)
	// if err != nil {
	// 	log.Fatalf("Failed to open raw socket: %v", err)
	// }
	// defer rs.Close()

	// // bind the socket to the network interface
	// rs.Bind(*networkInterface)
	// if err != nil {
	// 	log.Fatalf("Failed to bind socket: %v", err)
	// }
	// log.Printf("listening for %s packets on interface: %s\n", protocol, networkInterface.Name)

	// // SYS_RECVFROM syscall
	// rs.ReadPackets()
}
