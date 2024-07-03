package tui

import (
	"fmt"
	"net"
	"syscall"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ifaceItem represents a network interface in a list
type ifaceItem struct {
	name  string
	flags string
}

func (i ifaceItem) Title() string { return i.name }

func (i ifaceItem) Description() string { return i.flags }

func (i ifaceItem) FilterValue() string { return i.name }

// ifaceItem represents a network protocol in a list
type protoItem struct {
	name    string
	ethType uint16
}

func (p protoItem) Title() string { return p.name }

func (p protoItem) Description() string { return fmt.Sprintf("Eth type 0x%X", p.ethType) }

func (p protoItem) FilterValue() string { return p.name }

func fetchInterfaces() tea.Cmd {
	return func() tea.Msg {
		ifaces, err := net.Interfaces()
		if err != nil {
			return errMsg{err: err}
		}
		return networkInterfacesMsg(ifaces)
	}
}

type startMenuStep uint

const (
	ifaceStep startMenuStep = iota
	protoStep
	receiveStep
)

type startMenuModel struct {
	step          startMenuStep
	ifaceList     list.Model
	protoList     list.Model
	chosenIface   *net.Interface
	chosenProto   string
	chosenEthType uint16
	err           error
}

func NewStartMenuModel() *startMenuModel {
	const listHeight = 50
	const listWidth = 50

	items := []list.Item{}

	titlesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99")).Blink(true).Bold(true)
	ifaceDelegate := list.NewDefaultDelegate()
	ifaceDelegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	ifaceList := list.New(items, ifaceDelegate, listWidth, listHeight)
	ifaceList.Title = "Select a Network Interface"
	ifaceList.Styles.Title = titlesStyle
	ifaceList.SetShowStatusBar(true)
	ifaceList.SetFilteringEnabled(false)
	ifaceList.SetShowHelp(true)

	protoDelegate := list.NewDefaultDelegate()
	protoDelegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	protoList := list.New(items, protoDelegate, listWidth, listHeight)
	protoList.Title = "Select a Network Protocol"
	protoList.Styles.Title = titlesStyle
	protoList.SetShowStatusBar(false)
	protoList.SetFilteringEnabled(false)
	protoList.SetShowHelp(true)

	protoList.SetItems([]list.Item{
		protoItem{name: "all", ethType: syscall.ETH_P_ALL},
		protoItem{name: "arp", ethType: syscall.ETH_P_ARP},
		protoItem{name: "ip", ethType: syscall.ETH_P_IP},
		protoItem{name: "ipv6", ethType: syscall.ETH_P_IPV6},
		// UDP and TCP are part of IP, need special handling if filtered specifically
		protoItem{name: "udp", ethType: syscall.ETH_P_IP},
		protoItem{name: "udp6", ethType: syscall.ETH_P_IPV6},
		protoItem{name: "tcp", ethType: syscall.ETH_P_IP},
		protoItem{name: "tcp6", ethType: syscall.ETH_P_IPV6},
	})

	m := &startMenuModel{
		ifaceList: ifaceList,
		protoList: protoList,
	}

	return m
}

func (m startMenuModel) Init() tea.Cmd {
	return fetchInterfaces()
}

func (m startMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case networkInterfacesMsg:
		items := make([]list.Item, len(msg))
		for i, iface := range msg {
			items[i] = ifaceItem{
				name:  iface.Name,
				flags: iface.Flags.String(),
			}
		}
		cmd := m.ifaceList.SetItems(items)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.step {
			case ifaceStep:
				item, ok := m.ifaceList.SelectedItem().(ifaceItem)
				if ok {
					// retrieve the network interface
					networkInterface, err := net.InterfaceByName(item.name)
					if err != nil {
						m.err = err
						return m, tea.Quit
					}
					m.chosenIface = networkInterface
					m.step = protoStep
				}
			case protoStep:
				item, ok := m.protoList.SelectedItem().(protoItem)
				if ok {
					m.chosenProto = item.name
					m.chosenEthType = item.ethType
					m.step = receiveStep
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.step {
	case ifaceStep:
		m.ifaceList, cmd = m.ifaceList.Update(msg)
	case protoStep:
		m.protoList, cmd = m.protoList.Update(msg)
	}
	return m, cmd
}

func (m startMenuModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n", m.err)
	}

	var s string
	switch m.step {
	case ifaceStep:
		s = m.ifaceList.View()
	case protoStep:
		s = m.protoList.View()
	case receiveStep:
		s = fmt.Sprintf("You chose %s - %s\n", m.chosenIface.Name, m.chosenProto)
	default:
		return "Unkown step"
	}
	return lipgloss.NewStyle().Padding(1).Render(s)
}
