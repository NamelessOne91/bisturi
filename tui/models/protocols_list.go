package tui

import (
	"fmt"
	"syscall"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// protoItem represents a network protocol in a list
type protoItem struct {
	name    string
	ethType uint16
}

func (p protoItem) Title() string { return p.name }

func (p protoItem) Description() string { return fmt.Sprintf("Eth type 0x%X", p.ethType) }

func (p protoItem) FilterValue() string { return p.name }

type selectedProtocolItemMsg protoItem

type protocolsListModel struct {
	l list.Model
}

func (m protocolsListModel) Init() tea.Cmd {
	return nil
}

func (m protocolsListModel) Update(msg tea.Msg) (protocolsListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			p, ok := m.l.SelectedItem().(protoItem)
			if ok {
				return m, func() tea.Msg {
					return selectedProtocolItemMsg(p)
				}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.l, cmd = m.l.Update(msg)
	return m, cmd
}

func (m protocolsListModel) View() string {
	return m.l.View()
}

func newProtocolsListModel(width, height int) protocolsListModel {
	protoDelegate := list.NewDefaultDelegate()
	protoDelegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	items := []list.Item{
		protoItem{name: "all", ethType: syscall.ETH_P_ALL},
		protoItem{name: "arp", ethType: syscall.ETH_P_ARP},
		protoItem{name: "ip", ethType: syscall.ETH_P_IP},
		protoItem{name: "ipv6", ethType: syscall.ETH_P_IPV6},
		// UDP and TCP are part of IP, need special handling if filtered specifically
		protoItem{name: "udp", ethType: syscall.ETH_P_IP},
		protoItem{name: "udp6", ethType: syscall.ETH_P_IPV6},
		protoItem{name: "tcp", ethType: syscall.ETH_P_IP},
		protoItem{name: "tcp6", ethType: syscall.ETH_P_IPV6},
	}
	protoList := list.New(items, protoDelegate, width, height)
	protoList.Title = "Select a Network Protocol"
	protoList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99")).Blink(true).Bold(true)
	protoList.SetShowStatusBar(false)
	protoList.SetFilteringEnabled(false)
	protoList.SetShowHelp(true)

	return protocolsListModel{l: protoList}
}
