package main

import (
	"fmt"
	"log"
	"net"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type interfacesMenuModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m interfacesMenuModel) Init() tea.Cmd {
	return fetchInterfaces()
}

type networkInterfacesMsg []net.Interface

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

func fetchInterfaces() tea.Cmd {
	return func() tea.Msg {
		ifaces, err := net.Interfaces()
		if err != nil {
			return errMsg{err: err}
		}
		return networkInterfacesMsg(ifaces)
	}
}

type ifaceItem struct {
	name  string
	flags string
}

func (i ifaceItem) Title() string {
	return i.name
}

func (i ifaceItem) Description() string {
	return i.flags
}

func (i ifaceItem) FilterValue() string {
	return i.name
}

func (m interfacesMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.quitting = true
		return m, tea.Quit

	case networkInterfacesMsg:
		items := make([]list.Item, len(msg))
		for i, iface := range msg {
			items[i] = ifaceItem{
				name:  iface.Name,
				flags: iface.Flags.String(),
			}
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(ifaceItem)
			if ok {
				m.choice = i.name
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m interfacesMenuModel) View() string {
	if m.choice != "" {
		return fmt.Sprintf("You chose: %s\n", m.choice)
	}
	return lipgloss.NewStyle().Padding(1).Render(m.list.View())
}

func main() {
	const listHeight = 20
	const listWidth = 50

	items := []list.Item{}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New(items, delegate, listWidth, listHeight)
	l.Title = "Select a Network Interface"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)

	m := interfacesMenuModel{list: l}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}

	// retrieve the network interface
	// networkInterface, err := net.InterfaceByName(iface)
	// if err != nil {
	// 	log.Fatalf("Failed to get interface by name: %v", err)
	// }

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
